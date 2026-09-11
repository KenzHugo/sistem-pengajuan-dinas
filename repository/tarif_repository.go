package repository

import (
	"sipadin/models"

	"gorm.io/gorm"
)

// TarifRepository mengelola data tabel master_tarif_fasilitas
type TarifRepository interface {
	Create(tarif *models.MasterTarifFasilitas) error
	FindAll() ([]models.MasterTarifFasilitas, error)
	FindByID(id int) (*models.MasterTarifFasilitas, error)
	FindByKriteria(area, kriteria string) (*models.MasterTarifFasilitas, error)
	Update(tarif *models.MasterTarifFasilitas) error
	Delete(id int) error
}

type tarifRepository struct {
	db *gorm.DB
}

func NewTarifRepository(db *gorm.DB) TarifRepository {
	return &tarifRepository{db: db}
}

func (r *tarifRepository) Create(tarif *models.MasterTarifFasilitas) error {
	return r.db.Create(tarif).Error
}

func (r *tarifRepository) FindAll() ([]models.MasterTarifFasilitas, error) {
	var tarifs []models.MasterTarifFasilitas
	err := r.db.Order("area_tujuan ASC, kriteria_jabatan ASC").Find(&tarifs).Error
	return tarifs, err
}

func (r *tarifRepository) FindByID(id int) (*models.MasterTarifFasilitas, error) {
	var tarif models.MasterTarifFasilitas
	err := r.db.First(&tarif, id).Error
	if err != nil {
		return nil, err
	}
	return &tarif, nil
}

func (r *tarifRepository) FindByKriteria(area, kriteria string) (*models.MasterTarifFasilitas, error) {
	var tarif models.MasterTarifFasilitas
	err := r.db.Where("area_tujuan = ? AND kriteria_jabatan = ?", area, kriteria).First(&tarif).Error
	if err != nil {
		return nil, err
	}
	return &tarif, nil
}

func (r *tarifRepository) Update(tarif *models.MasterTarifFasilitas) error {
	return r.db.Save(tarif).Error
}

func (r *tarifRepository) Delete(id int) error {
	return r.db.Delete(&models.MasterTarifFasilitas{}, id).Error
}
