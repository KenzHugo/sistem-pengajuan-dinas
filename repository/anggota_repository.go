package repository

import (
	"sipadin/models"

	"gorm.io/gorm"
)

// AnggotaRepository mengelola data tabel sppd_anggota
type AnggotaRepository interface {
	Create(anggota *models.SppdAnggota) error
	FindBySppd(idSppd int) ([]models.SppdAnggota, error)
	DeleteBySppd(idSppd int) error
	Delete(id int) error
}

type anggotaRepository struct {
	db *gorm.DB
}

func NewAnggotaRepository(db *gorm.DB) AnggotaRepository {
	return &anggotaRepository{db: db}
}

func (r *anggotaRepository) Create(anggota *models.SppdAnggota) error {
	return r.db.Create(anggota).Error
}

func (r *anggotaRepository) FindBySppd(idSppd int) ([]models.SppdAnggota, error) {
	var anggota []models.SppdAnggota
	err := r.db.Preload("User").
		Where("id_sppd = ?", idSppd).
		Find(&anggota).Error
	return anggota, err
}

func (r *anggotaRepository) DeleteBySppd(idSppd int) error {
	return r.db.Where("id_sppd = ?", idSppd).Delete(&models.SppdAnggota{}).Error
}

func (r *anggotaRepository) Delete(id int) error {
	return r.db.Delete(&models.SppdAnggota{}, id).Error
}
