package service

import (
	"errors"
	"fmt"
	"net/url"
	"sipadin/dto"
	"sipadin/models"
	"sipadin/repository"
	"strings"
	"time"
)

// ApprovalService mendefinisikan business logic workflow approval 4 tahap
type ApprovalService interface {
	ManagerAction(idSppd, idApprover int, req dto.ApprovalRequest) (*models.SppdPengajuan, error)
	HRGAAction(idSppd, idApprover int, req dto.ApprovalRequest) (*models.SppdPengajuan, error)
	FinanceAction(idSppd, idApprover int, req dto.ApprovalRequest) (*models.SppdPengajuan, error)
	SendToDirector(idSppd, idApprover int) (*dto.SendWAResponse, error)
}

type approvalService struct {
	sppdRepo     repository.SppdRepository
	approvalRepo repository.ApprovalRepository
	userRepo     repository.UserRepository
}

func NewApprovalService(
	sppdRepo repository.SppdRepository,
	approvalRepo repository.ApprovalRepository,
	userRepo repository.UserRepository,
) ApprovalService {
	return &approvalService{
		sppdRepo:     sppdRepo,
		approvalRepo: approvalRepo,
		userRepo:     userRepo,
	}
}

// ManagerAction: pending_manager → pending_hrga (approve) atau rejected
func (s *approvalService) ManagerAction(idSppd, idApprover int, req dto.ApprovalRequest) (*models.SppdPengajuan, error) {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if sppd.StatusDokumen != "pending_manager" {
		return nil, errors.New("pengajuan tidak dalam status pending_manager")
	}

	var nextStatus, statusTindakan string
	if req.Action == "approve" {
		nextStatus = "pending_hrga"
		statusTindakan = "Approved"
	} else if req.Action == "reject" {
		nextStatus = "rejected"
		statusTindakan = "Rejected"
	} else {
		return nil, errors.New("action tidak valid, gunakan: approve atau reject")
	}

	return s.processApproval(idSppd, idApprover, "Manager", statusTindakan, nextStatus, req.Catatan)
}

// HRGAAction: pending_hrga → pending_finance (approve) atau rejected
func (s *approvalService) HRGAAction(idSppd, idApprover int, req dto.ApprovalRequest) (*models.SppdPengajuan, error) {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if sppd.StatusDokumen != "pending_hrga" {
		return nil, errors.New("pengajuan tidak dalam status pending_hrga")
	}

	var nextStatus, statusTindakan string
	if req.Action == "approve" {
		nextStatus = "pending_finance"
		statusTindakan = "Verified"
	} else if req.Action == "reject" {
		nextStatus = "rejected"
		statusTindakan = "Rejected"
	} else {
		return nil, errors.New("action tidak valid, gunakan: approve atau reject")
	}

	return s.processApproval(idSppd, idApprover, "HRGA", statusTindakan, nextStatus, req.Catatan)
}

// FinanceAction: pending_finance → ready_to_send (approve) atau rejected
func (s *approvalService) FinanceAction(idSppd, idApprover int, req dto.ApprovalRequest) (*models.SppdPengajuan, error) {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if sppd.StatusDokumen != "pending_finance" {
		return nil, errors.New("pengajuan tidak dalam status pending_finance")
	}

	var nextStatus, statusTindakan string
	if req.Action == "approve" {
		nextStatus = "ready_to_send"
		statusTindakan = "Approved"
	} else if req.Action == "reject" {
		nextStatus = "rejected"
		statusTindakan = "Rejected"
	} else {
		return nil, errors.New("action tidak valid, gunakan: approve atau reject")
	}

	return s.processApproval(idSppd, idApprover, "Finance", statusTindakan, nextStatus, req.Catatan)
}

// SendToDirector: HRGA mengirim dokumen ke WhatsApp Direksi → approved
func (s *approvalService) SendToDirector(idSppd, idApprover int) (*dto.SendWAResponse, error) {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if sppd.StatusDokumen != "ready_to_send" {
		return nil, errors.New("pengajuan belum siap dikirim (status harus ready_to_send)")
	}

	// Cari semua Direksi dengan nomor WA
	directors, err := s.userRepo.FindByRole("Direksi")
	if err != nil || len(directors) == 0 {
		return nil, errors.New("tidak ditemukan data Direksi dengan nomor WhatsApp")
	}

	director := directors[0]
	phone := strings.TrimPrefix(director.NoWhatsapp, "+")

	// Susun pesan WhatsApp
	tglBerangkat := sppd.TglBerangkat.Format("02 January 2006")
	tglKembali := sppd.TglKembali.Format("02 January 2006")
	pemohon := ""
	if sppd.Pemohon != nil {
		pemohon = sppd.Pemohon.NamaLengkap
	}

	msg := fmt.Sprintf(
		"*SPPD - Pengajuan Perjalanan Dinas*\n\n"+
			"Nomor SPPD : %s\n"+
			"Pemohon    : %s\n"+
			"Tujuan     : %s, %s\n"+
			"Tanggal    : %s s/d %s (%d hari)\n"+
			"Transportasi: %s\n"+
			"Akomodasi  : %s\n"+
			"Est. Biaya : Rp %s\n\n"+
			"Dokumen ini telah melalui verifikasi Manager, HRGA, dan Finance.\n"+
			"Mohon konfirmasi penerimaan. Terima kasih.",
		sppd.NoSppd,
		pemohon,
		sppd.KotaTujuan,
		sppd.InstansiTujuan,
		tglBerangkat, tglKembali, sppd.DurasiHari,
		sppd.Transportasi,
		sppd.Akomodasi,
		formatRupiah(sppd.EstimasiBiaya),
	)

	waURL := fmt.Sprintf("https://wa.me/%s?text=%s", phone, url.QueryEscape(msg))

	// Catat log dan update status ke approved
	logEntry := &models.SppdApprovalLog{
		IDSppd:         idSppd,
		TahapApproval:  "Direksi",
		IDApprover:     idApprover,
		StatusTindakan: "Sent",
		Catatan:        "Dokumen dikirim ke WhatsApp Direksi: " + director.NamaLengkap,
		TglTindakan:    time.Now(),
	}
	_ = s.approvalRepo.Create(logEntry)
	_ = s.sppdRepo.UpdateStatus(idSppd, "approved")

	return &dto.SendWAResponse{
		WAUrl:       waURL,
		NoSppd:      sppd.NoSppd,
		NamaDireksi: director.NamaLengkap,
	}, nil
}

// processApproval adalah helper: catat log + update status
func (s *approvalService) processApproval(
	idSppd, idApprover int,
	tahap, statusTindakan, nextStatus, catatan string,
) (*models.SppdPengajuan, error) {
	logEntry := &models.SppdApprovalLog{
		IDSppd:         idSppd,
		TahapApproval:  tahap,
		IDApprover:     idApprover,
		StatusTindakan: statusTindakan,
		Catatan:        catatan,
		TglTindakan:    time.Now(),
	}
	if err := s.approvalRepo.Create(logEntry); err != nil {
		return nil, errors.New("gagal menyimpan log persetujuan")
	}
	if err := s.sppdRepo.UpdateStatus(idSppd, nextStatus); err != nil {
		return nil, errors.New("gagal mengupdate status pengajuan")
	}
	return s.sppdRepo.FindByID(idSppd)
}

func formatRupiah(amount float64) string {
	s := fmt.Sprintf("%.0f", amount)
	n := len(s)
	result := ""
	for i, c := range s {
		if i > 0 && (n-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}
	return result
}
