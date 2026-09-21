package service

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"sipadin/config"
	"sipadin/dto"
	"sipadin/models"
	"sipadin/repository"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ApprovalService mendefinisikan business logic workflow approval 4 tahap
type ApprovalService interface {
	ManagerAction(idSppd, idApprover int, req dto.ApprovalRequest) (*models.SppdPengajuan, error)
	HRGAAction(idSppd, idApprover int, req dto.ApprovalRequest) (*models.SppdPengajuan, error)
	FinanceAction(idSppd, idApprover int, req dto.ApprovalRequest) (*models.SppdPengajuan, error)
	SendToDirector(idSppd, idApprover int) (*dto.SendWAResponse, error)
	GetSppdLedger(idSppd int) (*dto.SppdLedgerResponse, error)
}

type approvalService struct {
	sppdRepo     repository.SppdRepository
	approvalRepo repository.ApprovalRepository
	userRepo     repository.UserRepository
	notifSvc     NotificationService
	sppdSvc      SppdService
}

func NewApprovalService(
	sppdRepo repository.SppdRepository,
	approvalRepo repository.ApprovalRepository,
	userRepo repository.UserRepository,
	notifSvc NotificationService,
	sppdSvc SppdService,
) ApprovalService {
	return &approvalService{
		sppdRepo:     sppdRepo,
		approvalRepo: approvalRepo,
		userRepo:     userRepo,
		notifSvc:     notifSvc,
		sppdSvc:      sppdSvc,
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
		s.applyCostRevision(sppd, req, "HRGA", idApprover)
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
		s.applyCostRevision(sppd, req, "Finance", idApprover)
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
	directors, _ := s.userRepo.FindByRole("Direksi")

	var director models.User
	var phone string
	if len(directors) > 0 {
		for _, d := range directors {
			cleanPhone := strings.TrimSpace(d.NoWhatsapp)
			cleanPhone = strings.TrimPrefix(cleanPhone, "+")
			if cleanPhone != "" {
				director = d
				phone = cleanPhone
				break
			}
		}
		if director.NamaLengkap == "" {
			director = directors[0]
		}
	} else {
		director = models.User{
			NamaLengkap: "Direksi",
			NoWhatsapp:  "62895377328880",
		}
	}

	if phone == "" {
		phone = "62895377328880"
	}
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	if strings.HasPrefix(phone, "08") {
		phone = "628" + phone[2:]
	} else if strings.HasPrefix(phone, "8") {
		phone = "628" + phone[1:]
	}

	// Susun pesan WhatsApp
	tglBerangkat := sppd.TglBerangkat.Format("02 January 2006")
	tglKembali := sppd.TglKembali.Format("02 January 2006")
	pemohon := ""
	if sppd.Pemohon != nil {
		pemohon = sppd.Pemohon.NamaLengkap
	}

	// Format rombongan / rekan jika ada
	rombonganText := ""
	if len(sppd.Anggota) > 0 {
		var names []string
		for _, a := range sppd.Anggota {
			if a.User != nil && a.User.NamaLengkap != "" {
				names = append(names, a.User.NamaLengkap)
			}
		}
		if len(names) > 0 {
			rombonganText = fmt.Sprintf("Rombongan  : %s\n", strings.Join(names, ", "))
		}
	}

	lampiranText := ""
	if sppd.LampiranFile != "" {
		lampiranText = fmt.Sprintf("Lampiran   : %s\n", sppd.LampiranFile)
	}

	maksudText := ""
	if sppd.MaksudTujuan != "" {
		maksudText = fmt.Sprintf("Maksud     : %s\n", sppd.MaksudTujuan)
	}

	msg := fmt.Sprintf(
		"*SPPD - Pengajuan Perjalanan Dinas*\n\n"+
			"Nomor SPPD : %s\n"+
			"Pemohon    : %s\n"+
			"%s"+
			"Tujuan     : %s, %s\n"+
			"%s"+
			"Tanggal    : %s s/d %s (%d hari)\n"+
			"Transportasi: %s\n"+
			"Akomodasi  : %s\n"+
			"Est. Biaya : Rp %s\n"+
			"%s\n"+
			"Dokumen ini telah melalui verifikasi Manager, HRGA, dan Finance.\n"+
			"Mohon konfirmasi penerimaan. Terima kasih.",
		sppd.NoSppd,
		pemohon,
		rombonganText,
		sppd.KotaTujuan,
		sppd.InstansiTujuan,
		maksudText,
		tglBerangkat, tglKembali, sppd.DurasiHari,
		sppd.Transportasi,
		sppd.Akomodasi,
		formatRupiah(sppd.EstimasiBiaya),
		lampiranText,
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

	// Picu notifikasi otomatis ke Pemohon dan seluruh Rekan Kerja / Rombongan
	if s.notifSvc != nil {
		s.notifSvc.TriggerSPPDNotification(idSppd, "SENT_TO_DIRECTOR", idApprover, "Dokumen resmi disetujui Direksi")
	}

	return &dto.SendWAResponse{
		WAUrl:       waURL,
		NoSppd:      sppd.NoSppd,
		NamaDireksi: director.NamaLengkap,
	}, nil
}

// processApproval adalah helper: catat log + update status + rekam snapshot ledger keuangan dalam DB Transaction
func (s *approvalService) processApproval(
	idSppd, idApprover int,
	tahap, statusTindakan, nextStatus, catatan string,
) (*models.SppdPengajuan, error) {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		logEntry := &models.SppdApprovalLog{
			IDSppd:         idSppd,
			TahapApproval:  tahap,
			IDApprover:     idApprover,
			StatusTindakan: statusTindakan,
			Catatan:        catatan,
			TglTindakan:    time.Now(),
		}
		if err := tx.Create(logEntry).Error; err != nil {
			return errors.New("gagal menyimpan log persetujuan: " + err.Error())
		}
		if err := tx.Model(&models.SppdPengajuan{}).Where("id_sppd = ?", idSppd).Update("status_dokumen", nextStatus).Error; err != nil {
			return errors.New("gagal mengupdate status pengajuan: " + err.Error())
		}

		// Jika disetujui / diverifikasi oleh Manager, HRGA, atau Finance, rekam immutable snapshot ke sppd_detail_biaya
		if statusTindakan == "Approved" || statusTindakan == "Verified" {
			if errSnap := s.recordFinancialSnapshot(tx, idSppd, catatan); errSnap != nil {
				fmt.Printf("[Audit Ledger] Warning: Gagal rekam snapshot detail biaya: %v\n", errSnap)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Picu notifikasi otomatis berjenjang ke Approver berikutnya, Pemohon, dan Rekan Kerja
	if s.notifSvc != nil {
		var statusTahap string
		if nextStatus == "rejected" {
			statusTahap = "REJECTED"
		} else if tahap == "Manager" {
			statusTahap = "MANAGER_APPROVED"
		} else if tahap == "HRGA" {
			statusTahap = "HRGA_APPROVED"
		} else if tahap == "Finance" {
			statusTahap = "FINANCE_APPROVED"
		} else {
			statusTahap = nextStatus
		}
		s.notifSvc.TriggerSPPDNotification(idSppd, statusTahap, idApprover, catatan)
	}

	return s.sppdRepo.FindByID(idSppd)
}

// recordFinancialSnapshot merekam data profil karyawan & nominal final per kepala ke tabel sppd_detail_biaya
func (s *approvalService) recordFinancialSnapshot(tx *gorm.DB, idSppd int, catatanKoreksi string) error {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil || sppd == nil {
		return errors.New("data SPPD tidak ditemukan")
	}

	var rekanList []string
	for _, a := range sppd.Anggota {
		if a.User != nil && a.User.NomorID != "" {
			rekanList = append(rekanList, a.User.NomorID)
		} else if a.User != nil && a.User.NamaLengkap != "" {
			rekanList = append(rekanList, a.User.NamaLengkap)
		}
	}

	nomorIDPemohon := ""
	namaPemohon := ""
	kodeJabPemohon := ""

	if sppd.Pemohon != nil {
		nomorIDPemohon = strings.TrimSpace(sppd.Pemohon.NomorID)
		namaPemohon = strings.TrimSpace(sppd.Pemohon.NamaLengkap)
		kodeJabPemohon = strings.TrimSpace(sppd.Pemohon.KodeJabatan)
	}
	if nomorIDPemohon == "" && sppd.IDPemohon > 0 && s.userRepo != nil {
		if u, err := s.userRepo.FindByID(sppd.IDPemohon); err == nil && u != nil {
			nomorIDPemohon = strings.TrimSpace(u.NomorID)
			namaPemohon = strings.TrimSpace(u.NamaLengkap)
			kodeJabPemohon = strings.TrimSpace(u.KodeJabatan)
		}
	}

	var rincianPerKepala []dto.RincianPerKepalaDTO
	if s.sppdSvc != nil {
		calc, errCalc := s.sppdSvc.HitungRincianBiayaSPPD(
			nomorIDPemohon,
			sppd.JarakKM,
			sppd.JamBerangkat,
			sppd.JamKembali,
			sppd.TglBerangkat.Format("2006-01-02"),
			sppd.TglKembali.Format("2006-01-02"),
			sppd.KotaTujuan,
			sppd.Transportasi,
			sppd.Akomodasi,
			rekanList,
		)
		if errCalc == nil && calc != nil && len(calc.RincianPerKepala) > 0 {
			rincianPerKepala = calc.RincianPerKepala
		}
	}

	// Fallback jika kalkulasi dinamis belum terisi
	if len(rincianPerKepala) == 0 {
		if namaPemohon == "" {
			namaPemohon = "Pemohon"
		}
		rincianPerKepala = append(rincianPerKepala, dto.RincianPerKepalaDTO{
			NomorID:        nomorIDPemohon,
			Nama:           namaPemohon,
			Role:           "Pemohon",
			KodeJabatan:    kodeJabPemohon,
			NamaJabatan:    kodeJabPemohon,
			GroupDinasLuar: "Operator - Karu",
			Makan:          sppd.BiayaMakan,
			Saku:           sppd.BiayaUangSaku,
			Hotel:          sppd.BiayaHotel,
			Subtotal:       sppd.BiayaMakan + sppd.BiayaUangSaku + sppd.BiayaHotel,
		})
	}

	tujuanSnapshot := strings.TrimSpace(sppd.KotaTujuan)
	if sppd.InstansiTujuan != "" {
		tujuanSnapshot = fmt.Sprintf("%s (%s)", sppd.KotaTujuan, sppd.InstansiTujuan)
	}

	var ledgerItems []models.SppdDetailBiaya
	for _, p := range rincianPerKepala {
		subtotal := p.Makan + p.Saku + p.Hotel
		role := p.Role
		if strings.EqualFold(role, "pemohon") {
			role = "Pemohon"
		} else {
			role = "Rekan Kerja"
		}

		grupDinas := p.GroupDinasLuar
		if grupDinas == "" {
			grupDinas = p.NamaJabatan
		}

		item := models.SppdDetailBiaya{
			IDPengajuan:         sppd.IDSppd,
			IDKaryawan:          p.NomorID,
			RolePersonil:        role,
			SnapshotNama:        p.Nama,
			SnapshotJabatan:     p.NamaJabatan,
			SnapshotGrupDinas:   grupDinas,
			SnapshotTujuan:      tujuanSnapshot,
			SnapshotDeskripsi:   sppd.MaksudTujuan,
			UangMakan:           p.Makan,
			UangSaku:            p.Saku,
			UangHotel:           p.Hotel,
			SubtotalHakIndividu: subtotal,
			CatatanKoreksi:      catatanKoreksi,
		}
		ledgerItems = append(ledgerItems, item)
	}

	return s.sppdRepo.SaveDetailBiayaSnapshot(tx, sppd.IDSppd, ledgerItems)
}

// GetSppdLedger mengambil snapshot ledger keuangan komprehensif untuk halaman audit Finance
func (s *approvalService) GetSppdLedger(idSppd int) (*dto.SppdLedgerResponse, error) {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil || sppd == nil {
		return nil, errors.New("pengajuan SPPD tidak ditemukan")
	}

	// 1. Ambil detail biaya dengan WHERE id_pengajuan = ?
	details, err := s.sppdRepo.FindDetailBiayaByPengajuan(idSppd)
	if err != nil {
		return nil, err
	}

	// 2. Validasi integritas snapshot: jika belum ada atau snapshot lama corrupt/mismatch dengan pemohon SPPD
	needsRefresh := (len(details) == 0)
	if !needsRefresh && len(details) > 0 {
		var expectedPemohonID string
		if sppd.Pemohon != nil {
			expectedPemohonID = strings.TrimSpace(sppd.Pemohon.NomorID)
		} else if sppd.IDPemohon > 0 && s.userRepo != nil {
			if u, err := s.userRepo.FindByID(sppd.IDPemohon); err == nil && u != nil {
				expectedPemohonID = strings.TrimSpace(u.NomorID)
			}
		}
		if expectedPemohonID != "" && details[0].RolePersonil == "Pemohon" && strings.TrimSpace(details[0].IDKaryawan) != expectedPemohonID {
			needsRefresh = true
		}
	}

	if needsRefresh {
		_ = s.recordFinancialSnapshot(nil, idSppd, "")
		details, _ = s.sppdRepo.FindDetailBiayaByPengajuan(idSppd)
	}

	// 3. Rekalkulasi Grand Total di Backend: (Total Uang Makan + Saku + Hotel dari semua personil) + Biaya Transport/Tiket Bersama
	var totalHak float64
	for _, d := range details {
		totalHak += d.SubtotalHakIndividu
	}

	// Biaya Bersama ditarik dari kolom SPPD master yang berelasi langsung dengan ID ini
	biayaBersama := sppd.BiayaTransport + sppd.BiayaTiket
	grandTotalRekalkulasi := totalHak + biayaBersama

	// Sinkronisasi dengan estimasi_biaya master SPPD jika terdapat selisih
	if math.Abs(sppd.EstimasiBiaya-grandTotalRekalkulasi) > 1.0 && grandTotalRekalkulasi > 0 {
		sppd.EstimasiBiaya = grandTotalRekalkulasi
		_ = s.sppdRepo.Update(sppd)
	}

	resp := &dto.SppdLedgerResponse{
		IDSppd:             sppd.IDSppd,
		NoSppd:             sppd.NoSppd,
		KotaTujuan:         sppd.KotaTujuan,
		InstansiTujuan:     sppd.InstansiTujuan,
		MaksudTujuan:       sppd.MaksudTujuan,
		TglBerangkat:       sppd.TglBerangkat.Format("2006-01-02"),
		TglKembali:         sppd.TglKembali.Format("2006-01-02"),
		DurasiHari:         sppd.DurasiHari,
		DurasiJam:          sppd.DurasiJam,
		Transportasi:       sppd.Transportasi,
		Akomodasi:          sppd.Akomodasi,
		BiayaTransport:     sppd.BiayaTransport,
		BiayaTiket:         sppd.BiayaTiket,
		EstimasiBiayaTotal: grandTotalRekalkulasi,
		StatusDokumen:      sppd.StatusDokumen,
		TotalHakIndividu:   totalHak,
		TotalBersama:       biayaBersama,
		DetailLedger:       details,
	}

	return resp, nil
}

func (s *approvalService) applyCostRevision(sppd *models.SppdPengajuan, req dto.ApprovalRequest, role string, idApprover int) bool {
	nominalLama := sppd.EstimasiBiaya
	revised := false
	if req.BiayaMakan != nil && *req.BiayaMakan >= 0 {
		sppd.BiayaMakan = *req.BiayaMakan
		revised = true
	}
	if req.BiayaTransport != nil && *req.BiayaTransport >= 0 {
		sppd.BiayaTransport = *req.BiayaTransport
		revised = true
	}
	if req.BiayaUangSaku != nil && *req.BiayaUangSaku >= 0 {
		sppd.BiayaUangSaku = *req.BiayaUangSaku
		revised = true
	}
	if req.BiayaTiket != nil && *req.BiayaTiket >= 0 {
		sppd.BiayaTiket = *req.BiayaTiket
		revised = true
	}
	if req.BiayaHotel != nil && *req.BiayaHotel >= 0 {
		sppd.BiayaHotel = *req.BiayaHotel
		revised = true
	}

	totalKomponen := sppd.BiayaMakan + sppd.BiayaTransport + sppd.BiayaUangSaku + sppd.BiayaTiket + sppd.BiayaHotel
	if totalKomponen > 0 {
		sppd.EstimasiBiaya = totalKomponen
	} else if req.EstimasiBiaya != nil && *req.EstimasiBiaya >= 0 {
		sppd.EstimasiBiaya = *req.EstimasiBiaya
		revised = true
	}

	nominalBaru := sppd.EstimasiBiaya

	if revised {
		exceeded := false
		plafonTotal := sppd.PlafonMakan + sppd.PlafonTransport + sppd.PlafonUangSaku + sppd.PlafonTiket + sppd.PlafonHotel
		if (sppd.PlafonMakan > 0 && sppd.BiayaMakan > sppd.PlafonMakan) ||
			(sppd.PlafonTransport > 0 && sppd.BiayaTransport > sppd.PlafonTransport) ||
			(sppd.PlafonUangSaku > 0 && sppd.BiayaUangSaku > sppd.PlafonUangSaku) ||
			(sppd.PlafonTiket > 0 && sppd.BiayaTiket > sppd.PlafonTiket) ||
			(sppd.PlafonHotel > 0 && sppd.BiayaHotel > sppd.PlafonHotel) ||
			(plafonTotal > 0 && sppd.EstimasiBiaya > plafonTotal) {
			exceeded = true
		}
		if exceeded {
			sppd.StatusPlafon = "MELEBIHI_PLAFON"
		} else {
			sppd.StatusPlafon = "SESUAI"
		}
		sppd.RincianBiaya = fmt.Sprintf(
			"Revisi oleh %s: Makan %s | Trans %s | Saku %s | Tiket %s | Hotel %s | Total: %s (Status: %s)",
			role, formatRupiah(sppd.BiayaMakan), formatRupiah(sppd.BiayaTransport), formatRupiah(sppd.BiayaUangSaku),
			formatRupiah(sppd.BiayaTiket), formatRupiah(sppd.BiayaHotel), formatRupiah(sppd.EstimasiBiaya), sppd.StatusPlafon,
		)
		_ = s.sppdRepo.Update(sppd)

		// 3. AUDIT TRAIL / RIWAYAT PERUBAHAN
		// Jika terdapat perbedaan (nominal_baru != nominal_lama), INSERT ke sppd_approval_log
		if nominalBaru != nominalLama {
			approverLabel := role
			if s.userRepo != nil && idApprover > 0 {
				if approverUser, err := s.userRepo.FindByID(idApprover); err == nil && approverUser != nil && approverUser.NamaLengkap != "" {
					approverLabel = fmt.Sprintf("%s (%s)", approverUser.NamaLengkap, role)
				}
			}

			keterangan := fmt.Sprintf(
				"Estimasi biaya diubah oleh %s dari Rp%s menjadi Rp%s.",
				approverLabel,
				formatRupiah(nominalLama),
				formatRupiah(nominalBaru),
			)

			auditLog := &models.SppdApprovalLog{
				IDSppd:         sppd.IDSppd,
				TahapApproval:  role,
				IDApprover:     idApprover,
				StatusTindakan: "Biaya Diubah",
				Catatan:        keterangan,
				TglTindakan:    time.Now(),
			}
			_ = s.approvalRepo.Create(auditLog)
		}
	}
	return revised
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
