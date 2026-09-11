package models

import "time"

// SppdPengajuan merepresentasikan tabel `sppd_pengajuan` — master dokumen SPPD.
// Kolom persis sesuai PRD v2.0. TIDAK ada kolom tambahan.
type SppdPengajuan struct {
	IDSppd         int       `gorm:"column:id_sppd;primaryKey;autoIncrement"              json:"id_sppd"`
	NoSppd         string    `gorm:"column:no_sppd;type:varchar(50)"                      json:"no_sppd"`
	IDPemohon      int       `gorm:"column:id_pemohon;type:int;not null"                  json:"id_pemohon"`
	KotaTujuan     string    `gorm:"column:kota_tujuan;type:varchar(100);not null"        json:"kota_tujuan"`
	InstansiTujuan string    `gorm:"column:instansi_tujuan;type:varchar(150);not null"    json:"instansi_tujuan"`
	MaksudTujuan   string    `gorm:"column:maksud_tujuan;type:text;not null"              json:"maksud_tujuan"`
	TglBerangkat   time.Time `gorm:"column:tgl_berangkat;type:date;not null"              json:"tgl_berangkat"`
	TglKembali     time.Time `gorm:"column:tgl_kembali;type:date;not null"                json:"tgl_kembali"`
	DurasiHari     int       `gorm:"column:durasi_hari"                                   json:"durasi_hari"`
	Transportasi   string    `gorm:"column:transportasi;type:varchar(50)"                 json:"transportasi"`
	Akomodasi      string    `gorm:"column:akomodasi;type:varchar(50)"                    json:"akomodasi"`
	EstimasiBiaya  float64   `gorm:"column:estimasi_biaya;type:decimal(20,4)"             json:"estimasi_biaya"`
	StatusDokumen  string    `gorm:"column:status_dokumen;type:varchar(50);default:draft" json:"status_dokumen"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"                     json:"created_at"`

	// Relasi (preload only, tidak menjadi kolom di DB)
	Pemohon      *User             `gorm:"foreignKey:IDPemohon;references:IDUser"  json:"pemohon,omitempty"`
	ApprovalLogs []SppdApprovalLog `gorm:"foreignKey:IDSppd;references:IDSppd"     json:"approval_logs,omitempty"`
	Anggota      []SppdAnggota     `gorm:"foreignKey:IDSppd;references:IDSppd"     json:"anggota,omitempty"`
}

func (SppdPengajuan) TableName() string {
	return "sppd_pengajuan"
}
