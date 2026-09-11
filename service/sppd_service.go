package service

import (
	"errors"
	"sipadin/dto"
	"sipadin/models"
	"sipadin/repository"
	"time"
)

// SppdService mendefinisikan business logic pengajuan SPPD
type SppdService interface {
	Create(idPemohon int, req dto.CreateSppdRequest) (*models.SppdPengajuan, error)
	Submit(idSppd, idPemohon int) (*models.SppdPengajuan, error)
	GetByID(id int) (*models.SppdPengajuan, error)
	GetMyList(idPemohon int, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error)
	GetAllList(filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error)
	GetByStatus(statuses []string, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error)
	Update(idSppd, idPemohon int, req dto.UpdateSppdRequest) (*models.SppdPengajuan, error)
	Delete(idSppd, idPemohon int) error
}

type sppdService struct {
	sppdRepo    repository.SppdRepository
	anggotaRepo repository.AnggotaRepository
}

func NewSppdService(sppdRepo repository.SppdRepository, anggotaRepo repository.AnggotaRepository) SppdService {
	return &sppdService{sppdRepo: sppdRepo, anggotaRepo: anggotaRepo}
}

// Create membuat pengajuan SPPD baru dengan status draft
func (s *sppdService) Create(idPemohon int, req dto.CreateSppdRequest) (*models.SppdPengajuan, error) {
	tglBerangkat, err := time.Parse("2006-01-02", req.TglBerangkat)
	if err != nil {
		return nil, errors.New("format tgl_berangkat tidak valid, gunakan YYYY-MM-DD")
	}
	tglKembali, err := time.Parse("2006-01-02", req.TglKembali)
	if err != nil {
		return nil, errors.New("format tgl_kembali tidak valid, gunakan YYYY-MM-DD")
	}
	if tglKembali.Before(tglBerangkat) {
		return nil, errors.New("tgl_kembali tidak boleh sebelum tgl_berangkat")
	}

	durasi := int(tglKembali.Sub(tglBerangkat).Hours()/24) + 1

	sppd := &models.SppdPengajuan{
		IDPemohon:      idPemohon,
		KotaTujuan:     req.KotaTujuan,
		InstansiTujuan: req.InstansiTujuan,
		MaksudTujuan:   req.MaksudTujuan,
		TglBerangkat:   tglBerangkat,
		TglKembali:     tglKembali,
		DurasiHari:     durasi,
		Transportasi:   req.Transportasi,
		Akomodasi:      req.Akomodasi,
		EstimasiBiaya:  req.EstimasiBiaya,
		StatusDokumen:  "draft",
	}

	if err := s.sppdRepo.Create(sppd); err != nil {
		return nil, errors.New("gagal menyimpan pengajuan")
	}

	// Simpan anggota jika ada
	for _, a := range req.Anggota {
		anggota := &models.SppdAnggota{
			IDSppd:            sppd.IDSppd,
			IDUser:            a.IDUser,
			KriteriaFasilitas: a.KriteriaFasilitas,
			UangMakan:         a.UangMakan,
			UangSaku:          a.UangSaku,
		}
		_ = s.anggotaRepo.Create(anggota)
	}

	return s.sppdRepo.FindByID(sppd.IDSppd)
}

// Submit mengubah status dari draft ke pending_manager dan generate no_sppd
func (s *sppdService) Submit(idSppd, idPemohon int) (*models.SppdPengajuan, error) {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if sppd.IDPemohon != idPemohon {
		return nil, errors.New("anda tidak memiliki akses ke pengajuan ini")
	}
	if sppd.StatusDokumen != "draft" {
		return nil, errors.New("hanya pengajuan berstatus draft yang bisa diajukan")
	}

	noSppd, err := s.sppdRepo.GenerateNoSppd()
	if err != nil {
		return nil, errors.New("gagal generate nomor SPPD")
	}

	sppd.NoSppd = noSppd
	sppd.StatusDokumen = "pending_manager"
	if err := s.sppdRepo.Update(sppd); err != nil {
		return nil, errors.New("gagal mengajukan SPPD")
	}

	return s.sppdRepo.FindByID(idSppd)
}

func (s *sppdService) GetByID(id int) (*models.SppdPengajuan, error) {
	sppd, err := s.sppdRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	return sppd, nil
}

func (s *sppdService) GetMyList(idPemohon int, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error) {
	return s.sppdRepo.FindByPemohon(idPemohon, filter)
}

func (s *sppdService) GetAllList(filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error) {
	return s.sppdRepo.FindAll(filter)
}

func (s *sppdService) GetByStatus(statuses []string, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error) {
	return s.sppdRepo.FindByStatus(statuses, filter)
}

// Update mengedit pengajuan yang masih berstatus draft
func (s *sppdService) Update(idSppd, idPemohon int, req dto.UpdateSppdRequest) (*models.SppdPengajuan, error) {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if sppd.IDPemohon != idPemohon {
		return nil, errors.New("anda tidak memiliki akses ke pengajuan ini")
	}
	if sppd.StatusDokumen != "draft" {
		return nil, errors.New("hanya pengajuan berstatus draft yang bisa diedit")
	}

	if req.KotaTujuan != "" {
		sppd.KotaTujuan = req.KotaTujuan
	}
	if req.InstansiTujuan != "" {
		sppd.InstansiTujuan = req.InstansiTujuan
	}
	if req.MaksudTujuan != "" {
		sppd.MaksudTujuan = req.MaksudTujuan
	}
	if req.TglBerangkat != "" {
		t, _ := time.Parse("2006-01-02", req.TglBerangkat)
		sppd.TglBerangkat = t
	}
	if req.TglKembali != "" {
		t, _ := time.Parse("2006-01-02", req.TglKembali)
		sppd.TglKembali = t
	}
	sppd.DurasiHari = int(sppd.TglKembali.Sub(sppd.TglBerangkat).Hours()/24) + 1
	if req.Transportasi != "" {
		sppd.Transportasi = req.Transportasi
	}
	if req.Akomodasi != "" {
		sppd.Akomodasi = req.Akomodasi
	}
	if req.EstimasiBiaya > 0 {
		sppd.EstimasiBiaya = req.EstimasiBiaya
	}

	if err := s.sppdRepo.Update(sppd); err != nil {
		return nil, errors.New("gagal mengupdate pengajuan")
	}

	// Sync ulang anggota jika dikirim
	if len(req.Anggota) > 0 {
		_ = s.anggotaRepo.DeleteBySppd(idSppd)
		for _, a := range req.Anggota {
			anggota := &models.SppdAnggota{
				IDSppd:            idSppd,
				IDUser:            a.IDUser,
				KriteriaFasilitas: a.KriteriaFasilitas,
				UangMakan:         a.UangMakan,
				UangSaku:          a.UangSaku,
			}
			_ = s.anggotaRepo.Create(anggota)
		}
	}

	return s.sppdRepo.FindByID(idSppd)
}

// Delete menghapus pengajuan yang masih berstatus draft
func (s *sppdService) Delete(idSppd, idPemohon int) error {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil {
		return errors.New("pengajuan tidak ditemukan")
	}
	if sppd.IDPemohon != idPemohon {
		return errors.New("anda tidak memiliki akses ke pengajuan ini")
	}
	if sppd.StatusDokumen != "draft" {
		return errors.New("hanya pengajuan berstatus draft yang bisa dihapus")
	}

	_ = s.anggotaRepo.DeleteBySppd(idSppd)
	return s.sppdRepo.Delete(idSppd)
}
