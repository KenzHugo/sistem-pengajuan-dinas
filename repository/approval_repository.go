package repository

import (
	"sipadin/models"

	"gorm.io/gorm"
)

// ApprovalRepository mengelola data tabel sppd_approval_log
type ApprovalRepository interface {
	Create(log *models.SppdApprovalLog) error
	FindBySppd(idSppd int) ([]models.SppdApprovalLog, error)
	FindLastBySppd(idSppd int) (*models.SppdApprovalLog, error)
}

type approvalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) ApprovalRepository {
	return &approvalRepository{db: db}
}

func (r *approvalRepository) Create(log *models.SppdApprovalLog) error {
	return r.db.Create(log).Error
}

func (r *approvalRepository) FindBySppd(idSppd int) ([]models.SppdApprovalLog, error) {
	var logs []models.SppdApprovalLog
	err := r.db.
		Preload("Approver").
		Where("id_sppd = ?", idSppd).
		Order("tgl_tindakan ASC").
		Find(&logs).Error
	return logs, err
}

func (r *approvalRepository) FindLastBySppd(idSppd int) (*models.SppdApprovalLog, error) {
	var log models.SppdApprovalLog
	err := r.db.
		Where("id_sppd = ?", idSppd).
		Order("id_log DESC").
		First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}
