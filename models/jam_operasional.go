package models

import "time"

// JamOperasional merepresentasikan tabel `jamOperasional` di database adi_perdin
type JamOperasional struct {
	Hari        string     `gorm:"column:hari;primaryKey" json:"hari"`
	JamMasuk    *time.Time `gorm:"column:jamMasuk" json:"jam_masuk"`
	JamSelesai  *time.Time `gorm:"column:jamSelesai" json:"jam_selesai"`
	StatusLibur *bool      `gorm:"column:status_libur" json:"status_libur"`
}

// TableName menentukan nama tabel database
func (JamOperasional) TableName() string {
	return "jamOperasional"
}

// SpesialAkses merepresentasikan tabel `SpesialAkses` di database adi_perdin
type SpesialAkses struct {
	IDUser       int    `gorm:"column:id_user;primaryKey" json:"id_user"`
	NikPendek    string `gorm:"column:nik_pendek" json:"nik_pendek"`
	NamaLengkap  string `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	PasswordHash string `gorm:"column:password_hash" json:"-"`
}

// TableName menentukan nama tabel database
func (SpesialAkses) TableName() string {
	return "SpesialAkses"
}

