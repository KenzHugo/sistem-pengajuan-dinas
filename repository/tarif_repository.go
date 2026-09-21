package repository

import (
	"sipadin/models"

	"gorm.io/gorm"
)

// TarifRepository mengelola data tabel perjalanandinas_tarif dan perjalanandinas_fasilitas
type TarifRepository interface {
	FindAll(area string) ([]models.PerjalananDinasTarif, error)
	FindFasilitas() ([]models.PerjalananDinasFasilitas, error)
	FindByKey(area, jabatan string) (*models.PerjalananDinasTarif, error)
	GetAreas() ([]string, error)
	CreateOrUpdate(tarif *models.PerjalananDinasTarif) error
	Delete(area, jabatan string) error

	// Kompatibilitas helper
	Create(tarif *models.MasterTarifFasilitas) error
	FindByID(id int) (*models.MasterTarifFasilitas, error)
	FindByKriteria(area, kriteria string) (*models.MasterTarifFasilitas, error)
	Update(tarif *models.MasterTarifFasilitas) error
}

type tarifRepository struct {
	db *gorm.DB
}

// NewTarifRepository membuat instance baru TarifRepository
func NewTarifRepository(db *gorm.DB) TarifRepository {
	return &tarifRepository{db: db}
}

// FindAll mengambil daftar seluruh tarif dari tabel perjalanandinas_tarif
func (r *tarifRepository) FindAll(area string) ([]models.PerjalananDinasTarif, error) {
	var tarifs []models.PerjalananDinasTarif
	query := r.db.Model(&models.PerjalananDinasTarif{})
	if area != "" && area != "all" {
		query = query.Where("Area = ?", area)
	}
	err := query.Order("Area ASC, NoUrut ASC, Jabatan ASC").Find(&tarifs).Error
	return tarifs, err
}

// FindFasilitas mengambil seluruh kriteria hak fasilitas dari tabel perjalanandinas_fasilitas
func (r *tarifRepository) FindFasilitas() ([]models.PerjalananDinasFasilitas, error) {
	var fas []models.PerjalananDinasFasilitas
	err := r.db.Model(&models.PerjalananDinasFasilitas{}).Order("Kriteria ASC").Find(&fas).Error
	return fas, err
}

// FindByKey mencari tarif berdasarkan kombinasi Area dan Jabatan
func (r *tarifRepository) FindByKey(area, jabatan string) (*models.PerjalananDinasTarif, error) {
	var tarif models.PerjalananDinasTarif
	err := r.db.Where("Area = ? AND Jabatan = ?", area, jabatan).First(&tarif).Error
	if err != nil {
		return nil, err
	}
	return &tarif, nil
}

// GetAreas mengambil daftar unik Area tujuan
func (r *tarifRepository) GetAreas() ([]string, error) {
	var areas []string
	err := r.db.Model(&models.PerjalananDinasTarif{}).Distinct("Area").Order("Area ASC").Pluck("Area", &areas).Error
	return areas, err
}

// CreateOrUpdate menyimpan atau memperbarui data tarif di perjalanandinas_tarif
func (r *tarifRepository) CreateOrUpdate(tarif *models.PerjalananDinasTarif) error {
	return r.db.Save(tarif).Error
}

// Delete menghapus baris tarif di perjalanandinas_tarif
func (r *tarifRepository) Delete(area, jabatan string) error {
	return r.db.Where("Area = ? AND Jabatan = ?", area, jabatan).Delete(&models.PerjalananDinasTarif{}).Error
}

// Implementasi metode kompatibilitas
func (r *tarifRepository) Create(tarif *models.MasterTarifFasilitas) error {
	return r.db.Save(tarif).Error
}

func (r *tarifRepository) FindByID(id int) (*models.MasterTarifFasilitas, error) {
	var tarif models.MasterTarifFasilitas
	err := r.db.First(&tarif).Error
	return &tarif, err
}

func (r *tarifRepository) FindByKriteria(area, kriteria string) (*models.MasterTarifFasilitas, error) {
	var tarif models.MasterTarifFasilitas
	err := r.db.Where("Area = ? AND Jabatan LIKE ?", area, "%"+kriteria+"%").First(&tarif).Error
	return &tarif, err
}

func (r *tarifRepository) Update(tarif *models.MasterTarifFasilitas) error {
	return r.db.Save(tarif).Error
}
