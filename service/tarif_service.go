package service

import (
	"errors"
	"sipadin/dto"
	"sipadin/models"
	"sipadin/repository"
)

// TarifService mendefinisikan business logic master tarif fasilitas perjalanan dinas
type TarifService interface {
	GetAll(area string) ([]models.PerjalananDinasTarif, error)
	GetFasilitas() ([]models.PerjalananDinasFasilitas, error)
	GetAreas() ([]string, error)
	GetByKey(area, jabatan string) (*models.PerjalananDinasTarif, error)
	CreateOrUpdate(req dto.TarifRequest) (*models.PerjalananDinasTarif, error)
	Delete(area, jabatan string) error

	// Kompatibilitas
	GetByID(id int) (*models.MasterTarifFasilitas, error)
	Create(req dto.CreateTarifRequest) (*models.MasterTarifFasilitas, error)
	Update(id int, req dto.UpdateTarifRequest) (*models.MasterTarifFasilitas, error)
}

type tarifService struct {
	tarifRepo repository.TarifRepository
}

// NewTarifService membuat instance baru TarifService
func NewTarifService(tarifRepo repository.TarifRepository) TarifService {
	return &tarifService{tarifRepo: tarifRepo}
}

func (s *tarifService) GetAll(area string) ([]models.PerjalananDinasTarif, error) {
	return s.tarifRepo.FindAll(area)
}

func (s *tarifService) GetFasilitas() ([]models.PerjalananDinasFasilitas, error) {
	return s.tarifRepo.FindFasilitas()
}

func (s *tarifService) GetAreas() ([]string, error) {
	return s.tarifRepo.GetAreas()
}

func (s *tarifService) GetByKey(area, jabatan string) (*models.PerjalananDinasTarif, error) {
	tarif, err := s.tarifRepo.FindByKey(area, jabatan)
	if err != nil {
		return nil, errors.New("data tarif tidak ditemukan")
	}
	return tarif, nil
}

func (s *tarifService) CreateOrUpdate(req dto.TarifRequest) (*models.PerjalananDinasTarif, error) {
	area := req.GetArea()
	jabatan := req.GetJabatan()
	if area == "" || jabatan == "" {
		return nil, errors.New("area dan jabatan wajib diisi")
	}

	mataUang := req.MataUang
	if mataUang == "" {
		if area == "Jawa" || area == "Jakarta / Luar Jawa" {
			mataUang = "IDR"
		} else {
			mataUang = "USD"
		}
	}

	noUrut := req.NoUrut
	if noUrut == "" {
		noUrut = "1.0000"
	}

	tarif := &models.PerjalananDinasTarif{
		NoUrut:    noUrut,
		Area:      area,
		Jabatan:   jabatan,
		MataUang:  mataUang,
		UangMakan: req.GetUangMakan(),
		UangSaku:  req.GetUangSaku(),
		Hotel:     req.GetHotel(),
		Transport: req.GetTransport(),
		Tiket:     req.Tiket,
	}

	if err := s.tarifRepo.CreateOrUpdate(tarif); err != nil {
		return nil, errors.New("gagal menyimpan data tarif ke database: " + err.Error())
	}
	return tarif, nil
}

func (s *tarifService) Delete(area, jabatan string) error {
	if _, err := s.tarifRepo.FindByKey(area, jabatan); err != nil {
		return errors.New("data tarif tidak ditemukan")
	}
	return s.tarifRepo.Delete(area, jabatan)
}

// Implementasi kompatibilitas
func (s *tarifService) GetByID(id int) (*models.MasterTarifFasilitas, error) {
	return s.tarifRepo.FindByID(id)
}

func (s *tarifService) Create(req dto.CreateTarifRequest) (*models.MasterTarifFasilitas, error) {
	return s.CreateOrUpdate(req)
}

func (s *tarifService) Update(id int, req dto.UpdateTarifRequest) (*models.MasterTarifFasilitas, error) {
	return s.CreateOrUpdate(req)
}
