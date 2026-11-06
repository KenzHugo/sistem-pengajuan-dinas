package repository

import (
	"sipadin/models"
	"strings"

	"gorm.io/gorm"
)

// OperasionalRepository mendefinisikan kontrak akses data untuk tabel jamOperasional dan SpesialAkses
type OperasionalRepository interface {
	GetByHari(hari string) (*models.JamOperasional, error)
	GetAll() ([]models.JamOperasional, error)
	IsSpesialAkses(identifier string) bool
	GetSpesialAkses(identifier string) (*models.SpesialAkses, error)
}

type operasionalRepository struct {
	db *gorm.DB
}

// NewOperasionalRepository membuat instance baru OperasionalRepository
func NewOperasionalRepository(db *gorm.DB) OperasionalRepository {
	return &operasionalRepository{db: db}
}

// GetByHari mengambil konfigurasi jam operasional berdasarkan nama hari (misal: "Senin", "Minggu")
func (r *operasionalRepository) GetByHari(hari string) (*models.JamOperasional, error) {
	var item models.JamOperasional
	clean := strings.TrimSpace(hari)
	err := r.db.Where("LOWER(LTRIM(RTRIM(hari))) = LOWER(?)", clean).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetAll mengambil seluruh konfigurasi jam operasional
func (r *operasionalRepository) GetAll() ([]models.JamOperasional, error) {
	var list []models.JamOperasional
	err := r.db.Find(&list).Error
	return list, err
}

// IsSpesialAkses memeriksa apakah identifier (nomor_id, nik_pendek, atau nama) terdaftar di tabel SpesialAkses
func (r *operasionalRepository) IsSpesialAkses(identifier string) bool {
	sa, err := r.GetSpesialAkses(identifier)
	return err == nil && sa != nil
}

// GetSpesialAkses mengambil record dari tabel SpesialAkses berdasarkan nik_pendek, nomor_id, atau nama
func (r *operasionalRepository) GetSpesialAkses(identifier string) (*models.SpesialAkses, error) {
	clean := strings.TrimSpace(identifier)
	if clean == "" {
		return nil, gorm.ErrRecordNotFound
	}

	var sa models.SpesialAkses
	err := r.db.Where("nik_pendek = ? OR nik_pendek LIKE ? OR LOWER(nama_lengkap) = LOWER(?)", clean, "%"+clean, clean).First(&sa).Error
	if err != nil {
		return nil, err
	}
	return &sa, nil
}
