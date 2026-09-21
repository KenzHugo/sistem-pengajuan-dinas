package service

import (
	"errors"
	"math"
	"strings"
	"time"

	"sipadin/dto"
	"sipadin/models"
	"sipadin/repository"
)

// KaryawanService mendefinisikan business logic untuk entitas Karyawan
type KaryawanService interface {
	GetAll(query dto.KaryawanQuery) (*dto.KaryawanListResponse, error)
	GetByID(nomorID string) (*dto.KaryawanResponse, error)
	Create(req dto.CreateKaryawanRequest) (*dto.KaryawanResponse, error)
	Update(nomorID string, req dto.UpdateKaryawanRequest) (*dto.KaryawanResponse, error)
	Delete(nomorID string) error
	GetDepartemenList() ([]string, error)
	GetStats() (*dto.KaryawanStatsResponse, error)
}

type karyawanService struct {
	repo repository.KaryawanRepository
}

// NewKaryawanService membuat instance baru KaryawanService
func NewKaryawanService(repo repository.KaryawanRepository) KaryawanService {
	return &karyawanService{repo: repo}
}

// GetAll mengambil data karyawan berpaginasi dengan filter
func (s *karyawanService) GetAll(query dto.KaryawanQuery) (*dto.KaryawanListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 15
	}

	list, total, err := s.repo.FindAll(query)
	if err != nil {
		return nil, err
	}

	resList := make([]dto.KaryawanResponse, len(list))
	for i, k := range list {
		resList[i] = mapToKaryawanResponse(&k)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	depts, _ := s.repo.GetDepartemenList()

	return &dto.KaryawanListResponse{
		Total:       total,
		Page:        query.Page,
		Limit:       query.Limit,
		TotalPages:  totalPages,
		Data:        resList,
		Departemens: depts,
	}, nil
}

// GetByID mengambil detail karyawan berdasarkan NomorID
func (s *karyawanService) GetByID(nomorID string) (*dto.KaryawanResponse, error) {
	nomorID = strings.TrimSpace(nomorID)
	if nomorID == "" {
		return nil, errors.New("Nomor ID tidak boleh kosong")
	}

	karyawan, err := s.repo.FindByID(nomorID)
	if err != nil {
		return nil, errors.New("Data karyawan tidak ditemukan")
	}

	res := mapToKaryawanResponse(karyawan)
	return &res, nil
}

// Create membuat data karyawan baru
func (s *karyawanService) Create(req dto.CreateKaryawanRequest) (*dto.KaryawanResponse, error) {
	nomorID := strings.TrimSpace(req.NomorID)
	nama := strings.TrimSpace(req.Nama)

	if nomorID == "" {
		return nil, errors.New("Nomor ID wajib diisi")
	}
	if nama == "" {
		return nil, errors.New("Nama karyawan wajib diisi")
	}

	// Cek apakah NomorID sudah terdaftar
	existing, _ := s.repo.FindByID(nomorID)
	if existing != nil {
		return nil, errors.New("Nomor ID sudah terdaftar dalam sistem")
	}

	tglLahir, err := parseTanggalLahir(req.TanggalLahir)
	if err != nil {
		return nil, err
	}

	aktifStatus := 1
	if req.Aktif != nil {
		aktifStatus = *req.Aktif
	}

	karyawan := models.Karyawan{
		NomorID:        nomorID,
		Nama:           nama,
		KodeDepartemen: strings.TrimSpace(req.KodeDepartemen),
		KodeJabatan:    strings.TrimSpace(req.KodeJabatan),
		Job:            strings.TrimSpace(req.Job),
		Kota:           strings.TrimSpace(req.Kota),
		TanggalLahir:   tglLahir,
		Aktif:          aktifStatus,
	}

	if err := s.repo.Create(&karyawan); err != nil {
		return nil, err
	}

	res := mapToKaryawanResponse(&karyawan)
	return &res, nil
}

// Update memperbarui data karyawan
func (s *karyawanService) Update(nomorID string, req dto.UpdateKaryawanRequest) (*dto.KaryawanResponse, error) {
	nomorID = strings.TrimSpace(nomorID)
	nama := strings.TrimSpace(req.Nama)

	if nomorID == "" {
		return nil, errors.New("Nomor ID tidak boleh kosong")
	}
	if nama == "" {
		return nil, errors.New("Nama karyawan tidak boleh kosong")
	}

	existing, err := s.repo.FindByID(nomorID)
	if err != nil || existing == nil {
		return nil, errors.New("Data karyawan tidak ditemukan")
	}

	tglLahir, err := parseTanggalLahir(req.TanggalLahir)
	if err != nil {
		return nil, err
	}

	aktifStatus := existing.Aktif
	if req.Aktif != nil {
		aktifStatus = *req.Aktif
	}

	updated := models.Karyawan{
		NomorID:        existing.NomorID,
		Nama:           nama,
		KodeDepartemen: strings.TrimSpace(req.KodeDepartemen),
		KodeJabatan:    strings.TrimSpace(req.KodeJabatan),
		Job:            strings.TrimSpace(req.Job),
		Kota:           strings.TrimSpace(req.Kota),
		TanggalLahir:   tglLahir,
		Aktif:          aktifStatus,
	}

	if err := s.repo.Update(nomorID, &updated); err != nil {
		return nil, err
	}

	res := mapToKaryawanResponse(&updated)
	return &res, nil
}

// Delete menghapus data karyawan
func (s *karyawanService) Delete(nomorID string) error {
	nomorID = strings.TrimSpace(nomorID)
	if nomorID == "" {
		return errors.New("Nomor ID tidak boleh kosong")
	}

	existing, err := s.repo.FindByID(nomorID)
	if err != nil || existing == nil {
		return errors.New("Data karyawan tidak ditemukan")
	}

	return s.repo.Delete(nomorID)
}

// GetDepartemenList mengambil daftar unik kode departemen
func (s *karyawanService) GetDepartemenList() ([]string, error) {
	return s.repo.GetDepartemenList()
}

// GetStats mengambil ringkasan metrik karyawan
func (s *karyawanService) GetStats() (*dto.KaryawanStatsResponse, error) {
	return s.repo.GetStats()
}

// parseTanggalLahir mem-parse string tanggal dengan aman (null-safe)
func parseTanggalLahir(tglStr string) (*time.Time, error) {
	clean := strings.TrimSpace(tglStr)
	if clean == "" || clean == "null" || clean == "0000-00-00" {
		return nil, nil
	}

	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"02/01/2006",
		"02-01-2006",
	}

	for _, fmtStr := range formats {
		if t, err := time.Parse(fmtStr, clean); err == nil {
			return &t, nil
		}
	}

	return nil, errors.New("Format tanggal lahir tidak valid (gunakan format YYYY-MM-DD)")
}

// mapToKaryawanResponse memetakan model Karyawan ke DTO response
func mapToKaryawanResponse(k *models.Karyawan) dto.KaryawanResponse {
	var tglStr *string
	var tglFormatted string = "-"

	if k.TanggalLahir != nil && !k.TanggalLahir.IsZero() {
		iso := k.TanggalLahir.Format("2006-01-02")
		tglStr = &iso
		tglFormatted = k.TanggalLahir.Format("02 Jan 2006")
	}

	return dto.KaryawanResponse{
		NomorID:              k.NomorID,
		Nama:                 k.Nama,
		KodeDepartemen:       k.KodeDepartemen,
		KodeJabatan:          k.KodeJabatan,
		Job:                  k.Job,
		Kota:                 k.Kota,
		TanggalLahir:         tglStr,
		TanggalLahirFormatted: tglFormatted,
		Aktif:                k.Aktif,
	}
}
