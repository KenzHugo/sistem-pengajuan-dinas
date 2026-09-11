package repository

import (
	"fmt"
	"sipadin/dto"
	"sipadin/models"
	"time"

	"gorm.io/gorm"
)

// SppdRepository mendefinisikan kontrak akses data untuk tabel sppd_pengajuan
type SppdRepository interface {
	Create(sppd *models.SppdPengajuan) error
	FindByID(id int) (*models.SppdPengajuan, error)
	FindAll(filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error)
	FindByPemohon(idPemohon int, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error)
	FindByStatus(statuses []string, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error)
	Update(sppd *models.SppdPengajuan) error
	UpdateStatus(id int, status string) error
	Delete(id int) error
	GenerateNoSppd() (string, error)
}

type sppdRepository struct {
	db *gorm.DB
}

func NewSppdRepository(db *gorm.DB) SppdRepository {
	return &sppdRepository{db: db}
}

func (r *sppdRepository) Create(sppd *models.SppdPengajuan) error {
	return r.db.Create(sppd).Error
}

func (r *sppdRepository) FindByID(id int) (*models.SppdPengajuan, error) {
	var sppd models.SppdPengajuan
	err := r.db.
		Preload("Pemohon").
		Preload("ApprovalLogs").
		Preload("ApprovalLogs.Approver").
		Preload("Anggota").
		Preload("Anggota.User").
		First(&sppd, id).Error
	if err != nil {
		return nil, err
	}
	return &sppd, nil
}

func (r *sppdRepository) FindAll(filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error) {
	var sppds []models.SppdPengajuan
	var total int64

	query := r.db.Model(&models.SppdPengajuan{}).Preload("Pemohon")
	if filter.Status != "" {
		query = query.Where("status_dokumen = ?", filter.Status)
	}

	query.Count(&total)

	page, limit := parsePagination(filter.Page, filter.Limit)
	err := query.Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&sppds).Error

	return sppds, total, err
}

func (r *sppdRepository) FindByPemohon(idPemohon int, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error) {
	var sppds []models.SppdPengajuan
	var total int64

	query := r.db.Model(&models.SppdPengajuan{}).
		Preload("Pemohon").
		Where("id_pemohon = ?", idPemohon)
	if filter.Status != "" {
		query = query.Where("status_dokumen = ?", filter.Status)
	}

	query.Count(&total)

	page, limit := parsePagination(filter.Page, filter.Limit)
	err := query.Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&sppds).Error

	return sppds, total, err
}

func (r *sppdRepository) FindByStatus(statuses []string, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error) {
	var sppds []models.SppdPengajuan
	var total int64

	query := r.db.Model(&models.SppdPengajuan{}).
		Preload("Pemohon").
		Where("status_dokumen IN ?", statuses)

	query.Count(&total)

	page, limit := parsePagination(filter.Page, filter.Limit)
	err := query.Order("created_at ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&sppds).Error

	return sppds, total, err
}

func (r *sppdRepository) Update(sppd *models.SppdPengajuan) error {
	return r.db.Save(sppd).Error
}

func (r *sppdRepository) UpdateStatus(id int, status string) error {
	return r.db.Model(&models.SppdPengajuan{}).
		Where("id_sppd = ?", id).
		Update("status_dokumen", status).Error
}

func (r *sppdRepository) Delete(id int) error {
	return r.db.Delete(&models.SppdPengajuan{}, id).Error
}

// GenerateNoSppd membuat nomor SPPD dengan format SPPD/YYYYMM/NNNN
func (r *sppdRepository) GenerateNoSppd() (string, error) {
	now := time.Now()
	yearMonth := now.Format("200601")
	pattern := fmt.Sprintf("SPPD/%s/%%", yearMonth)

	var count int64
	r.db.Model(&models.SppdPengajuan{}).
		Where("no_sppd LIKE ?", pattern).
		Count(&count)

	seq := count + 1
	noSppd := fmt.Sprintf("SPPD/%s/%04d", yearMonth, seq)
	return noSppd, nil
}

func parsePagination(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return page, limit
}
