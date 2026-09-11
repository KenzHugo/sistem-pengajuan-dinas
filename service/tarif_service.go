package service

import (
	"errors"
	"sipadin/dto"
	"sipadin/models"
	"sipadin/repository"
)

// TarifService mendefinisikan business logic master tarif fasilitas
type TarifService interface {
	GetAll() ([]models.MasterTarifFasilitas, error)
	GetByID(id int) (*models.MasterTarifFasilitas, error)
	Create(req dto.CreateTarifRequest) (*models.MasterTarifFasilitas, error)
	Update(id int, req dto.UpdateTarifRequest) (*models.MasterTarifFasilitas, error)
	Delete(id int) error
}

type tarifService struct {
	tarifRepo repository.TarifRepository
}

func NewTarifService(tarifRepo repository.TarifRepository) TarifService {
	return &tarifService{tarifRepo: tarifRepo}
}

func (s *tarifService) GetAll() ([]models.MasterTarifFasilitas, error) {
	return s.tarifRepo.FindAll()
}

func (s *tarifService) GetByID(id int) (*models.MasterTarifFasilitas, error) {
	tarif, err := s.tarifRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("data tarif tidak ditemukan")
	}
	return tarif, nil
}

func (s *tarifService) Create(req dto.CreateTarifRequest) (*models.MasterTarifFasilitas, error) {
	tarif := &models.MasterTarifFasilitas{
		AreaTujuan:      req.AreaTujuan,
		KriteriaJabatan: req.KriteriaJabatan,
		HakTransport:    req.HakTransport,
		HakHotel:        req.HakHotel,
		TarifUangMakan:  req.TarifUangMakan,
		TarifUangSaku:   req.TarifUangSaku,
	}
	if err := s.tarifRepo.Create(tarif); err != nil {
		return nil, errors.New("gagal menyimpan data tarif")
	}
	return tarif, nil
}

func (s *tarifService) Update(id int, req dto.UpdateTarifRequest) (*models.MasterTarifFasilitas, error) {
	tarif, err := s.tarifRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("data tarif tidak ditemukan")
	}

	if req.AreaTujuan != "" {
		tarif.AreaTujuan = req.AreaTujuan
	}
	if req.KriteriaJabatan != "" {
		tarif.KriteriaJabatan = req.KriteriaJabatan
	}
	if req.HakTransport > 0 {
		tarif.HakTransport = req.HakTransport
	}
	if req.HakHotel > 0 {
		tarif.HakHotel = req.HakHotel
	}
	if req.TarifUangMakan > 0 {
		tarif.TarifUangMakan = req.TarifUangMakan
	}
	if req.TarifUangSaku > 0 {
		tarif.TarifUangSaku = req.TarifUangSaku
	}

	if err := s.tarifRepo.Update(tarif); err != nil {
		return nil, errors.New("gagal mengupdate data tarif")
	}
	return tarif, nil
}

func (s *tarifService) Delete(id int) error {
	if _, err := s.tarifRepo.FindByID(id); err != nil {
		return errors.New("data tarif tidak ditemukan")
	}
	return s.tarifRepo.Delete(id)
}
