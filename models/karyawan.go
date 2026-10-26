package models

import (
	"time"
)

// Karyawan merepresentasikan tabel `karyawan` di database adi_persons SQL Server
// Menggunakan pointer *time.Time untuk TanggalLahir agar null-safe saat berhadapan dengan nilai NULL di SQL Server
type Karyawan struct {
	NomorID        string     `gorm:"column:NomorID;primaryKey;type:varchar(50)" json:"nomor_id"`
	Nama           string     `gorm:"column:Nama;type:varchar(50)" json:"nama"`
	KodeDepartemen string     `gorm:"column:KodeDepartemen;type:varchar(16)" json:"kode_departemen"`
	KodeJabatan    string     `gorm:"column:KodeJabatan;type:varchar(16)" json:"kode_jabatan"`
	Job            string     `gorm:"column:Job;type:varchar(50)" json:"job"`
	Kota           string     `gorm:"column:Kota;type:varchar(32)" json:"kota"`
	TanggalLahir   *time.Time `gorm:"column:TanggalLahir;type:datetime" json:"tanggal_lahir"`
	Aktif          int        `gorm:"column:Aktif;type:int;default:1" json:"aktif"`
	Password       string     `gorm:"column:Password;type:varchar(255)" json:"-"`
}

// TableName menentukan nama tabel database secara eksplisit
func (Karyawan) TableName() string {
	// return "karyawan_copy1"
	return "adi_persons.dbo.karyawan"
}
