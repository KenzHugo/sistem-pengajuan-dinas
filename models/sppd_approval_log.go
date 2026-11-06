package models

import "time"

// SppdApprovalLog merepresentasikan tabel `sppd_approval_log` — riwayat persetujuan.
// Kolom persis sesuai PRD v2.0. TIDAK ada kolom tambahan.
type SppdApprovalLog struct {
	IDLog          int       `gorm:"column:id_log;primaryKey;autoIncrement"            json:"id_log"`
	IDSppd         int       `gorm:"column:id_sppd;type:int;not null"                          json:"id_sppd"`
	TahapApproval  string    `gorm:"column:tahap_approval;type:varchar(50);not null"   json:"tahap_approval"`
	IDApprover     int       `gorm:"column:id_approver;type:int;not null"                       json:"id_approver"`
	StatusTindakan string    `gorm:"column:status_tindakan;type:varchar(20);not null"  json:"status_tindakan"`
	Catatan        string    `gorm:"column:catatan;type:varchar(max)"                  json:"catatan"`
	TglTindakan    time.Time `gorm:"column:tgl_tindakan;autoCreateTime"                json:"tgl_tindakan"`

	// Relasi
	Approver *User `gorm:"foreignKey:IDApprover;references:IDUser" json:"approver,omitempty"`
}

func (SppdApprovalLog) TableName() string {
	return "sppd_approval_log"
}
