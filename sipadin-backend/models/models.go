package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Nama       string    `gorm:"size:100;not null" json:"nama"`
	NIP        string    `gorm:"size:50;unique;not null" json:"nip"`
	Email      string    `gorm:"size:100;unique;not null" json:"email"`
	Password   string    `gorm:"size:255;not null" json:"-"` // tidak di-return ke JSON
	Role       string    `gorm:"type:enum('karyawan','atasan','hrga');not null" json:"role"`
	Departemen string    `gorm:"size:100;not null" json:"departemen"`
	Jabatan    string    `gorm:"size:100;not null" json:"jabatan"`
	CreatedAt  time.Time `json:"created_at"`
}

type PengajuanDinas struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"not null" json:"user_id"`
	User             User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tujuan           string    `gorm:"size:255;not null" json:"tujuan"`
	Instansi         string    `gorm:"size:255;not null" json:"instansi"`
	TglKeberangkatan string    `gorm:"type:date;not null" json:"tgl_keberangkatan"`
	TglKembali       string    `gorm:"type:date;not null" json:"tgl_kembali"`
	Transportasi     string    `gorm:"size:100;not null" json:"transportasi"`
	Akomodasi        string    `gorm:"size:100;not null" json:"akomodasi"`
	Keperluan        string    `gorm:"type:text;not null" json:"keperluan"`
	Keterangan       string    `gorm:"type:text" json:"keterangan"`
	EstimasiBiaya    float64   `gorm:"type:decimal(12,2);not null" json:"estimasi_biaya"`
	Status           string    `gorm:"type:enum('pending','approved','rejected','verified_hrga','rejected_hrga','notified_direksi');default:'pending'" json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type NotifikasiDireksi struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	PengajuanID      uint      `gorm:"not null" json:"pengajuan_id"`
	IsiPesan         string    `gorm:"type:text;not null" json:"isi_pesan"`
	StatusPengiriman string    `gorm:"size:50;default:'terkirim'" json:"status_pengiriman"`
	SentAt           time.Time `json:"sent_at"`
}