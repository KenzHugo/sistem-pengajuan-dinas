package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey;column:id" json:"id"`
	Nama       string    `gorm:"column:nama;size:100;not null" json:"nama"`
	NIP        string    `gorm:"column:nip;size:50;not null" json:"nip"`
	Email      string    `gorm:"column:email;size:100;unique;not null" json:"email"`
	Password   string    `gorm:"column:password;size:255;not null" json:"password,omitempty"`
	Role       string    `gorm:"column:role;type:varchar(50);not null" json:"role"`
	Departemen string    `gorm:"column:departemen;size:100;not null" json:"departemen"`
	Jabatan    string    `gorm:"column:jabatan;size:100;not null" json:"jabatan"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func (User) TableName() string {
	return "users"
}

type PengajuanDinas struct {
	ID                     uint      `gorm:"primaryKey;column:id" json:"id"`
	NomorSurat             string    `gorm:"column:nomor_surat;size:100" json:"nomor_surat"`
	UserID                 uint      `gorm:"column:user_id;not null" json:"user_id"`
	User                   User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Nama                   string    `gorm:"column:nama;size:100" json:"nama"`
	NIP                    string    `gorm:"column:nip;size:50" json:"nip"`
	Dept                   string    `gorm:"column:dept;size:100" json:"dept"`
	Jabatan                string    `gorm:"column:jabatan;size:100" json:"jabatan"`
	Asal                   string    `gorm:"column:asal;size:100" json:"asal"`
	Tujuan                 string    `gorm:"column:tujuan;size:255;not null" json:"tujuan"`
	Instansi               string    `gorm:"column:instansi;size:255;not null" json:"instansi"`
	TglKeberangkatan       string    `gorm:"column:tgl_keberangkatan;size:50;not null" json:"tgl_keberangkatan"`
	TglKembali             string    `gorm:"column:tgl_kembali;size:50;not null" json:"tgl_kembali"`
	JamBerangkat           string    `gorm:"column:jam_berangkat;size:20" json:"jam_berangkat"`
	Jarak                  int       `gorm:"column:jarak" json:"jarak"`
	LamaJam                int       `gorm:"column:lama_jam" json:"lama_jam"`
	Durasi                 int       `gorm:"column:durasi" json:"durasi"`
	StatusMenginap         string    `gorm:"column:status_menginap;size:20" json:"status_menginap"`
	Kriteria               string    `gorm:"column:kriteria;size:100" json:"kriteria"`
	Transportasi           string    `gorm:"column:transportasi;size:100;not null" json:"transportasi"`
	Akomodasi              string    `gorm:"column:akomodasi;size:100" json:"akomodasi"`
	Keperluan              string    `gorm:"column:keperluan;type:text;not null" json:"keperluan"`
	Keterangan             string    `gorm:"column:keterangan;type:text" json:"keterangan"`
	BiayaTransport         float64   `gorm:"column:biaya_transport;type:decimal(12,2)" json:"biaya_transport"`
	BiayaPenginapan        float64   `gorm:"column:biaya_penginapan;type:decimal(12,2)" json:"biaya_penginapan"`
	UangHarian             float64   `gorm:"column:uang_harian;type:decimal(12,2)" json:"uang_harian"`
	EstimasiBiaya          float64   `gorm:"column:estimasi_biaya;type:decimal(12,2);not null" json:"estimasi_biaya"`
	Lampiran               string    `gorm:"column:lampiran;size:255" json:"lampiran"`
	Status                 string    `gorm:"column:status;type:varchar(50);default:'pending'" json:"status"`
	TanggalPengajuan       string    `gorm:"column:tanggal_pengajuan;size:50" json:"tanggal_pengajuan"`
	CatatanAtasan          string    `gorm:"column:catatan_atasan;type:text" json:"catatan_atasan"`
	DisetujuiOleh          string    `gorm:"column:disetujui_oleh;size:100" json:"disetujui_oleh"`
	TanggalDisetujui       string    `gorm:"column:tanggal_disetujui;size:50" json:"tanggal_disetujui"`
	CatatanHRGA            string    `gorm:"column:catatan_hrga;type:text" json:"catatan_hrga"`
	DiverifikasiOleh       string    `gorm:"column:diverifikasi_oleh;size:100" json:"diverifikasi_oleh"`
	TanggalVerifikasiHRGA  string    `gorm:"column:tanggal_verifikasi_hrga;size:50" json:"tanggal_verifikasi_hrga"`
	TanggalNotifikasiDireksi string  `gorm:"column:tanggal_notifikasi_direksi;size:50" json:"tanggal_notifikasi_direksi"`
	PesanWA                string    `gorm:"column:pesan_wa;type:text" json:"pesan_wa"`
	CreatedAt              time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (PengajuanDinas) TableName() string {
	return "pengajuan_dinas"
}

type NotifikasiDireksi struct {
	ID               uint      `gorm:"primaryKey;column:id" json:"id"`
	PengajuanID      uint      `gorm:"column:pengajuan_id;not null" json:"pengajuan_id"`
	IsiPesan         string    `gorm:"column:isi_pesan;type:text;not null" json:"isi_pesan"`
	StatusPengiriman string    `gorm:"column:status_pengiriman;size:50;default:'terkirim'" json:"status_pengiriman"`
	SentAt           time.Time `gorm:"column:sent_at" json:"sent_at"`
}

func (NotifikasiDireksi) TableName() string {
	return "notifikasi_direksi"
}