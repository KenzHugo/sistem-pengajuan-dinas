package models

// Jabatan merepresentasikan tabel master `jabatan` di MySQL
// Menjadi kunci relasi antara tabel `karyawan_copy1` (melalui KodeJabatan) dan `perjalanandinas_tarif` (melalui GroupDinasLuar)
type Jabatan struct {
	Kode           string `gorm:"column:Kode;primaryKey;type:varchar(10)" json:"kode"`
	Nama           string `gorm:"column:Nama;type:varchar(100)"          json:"nama"`
	GroupDinasLuar string `gorm:"column:GroupDinasLuar;type:varchar(50)" json:"group_dinas_luar"`
}

// TableName mendefinisikan nama tabel secara eksplisit
func (Jabatan) TableName() string {
	return "jabatan"
}
