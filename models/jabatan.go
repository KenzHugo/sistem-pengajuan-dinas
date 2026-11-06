package models

// Jabatan merepresentasikan tabel master `jabatan` di SQL Server (database adi_persons)
// Menjadi kunci relasi antara tabel `karyawan` (melalui KodeJabatan) dan `perjalanandinas_tarif` (melalui GroupDinasLuar)
type Jabatan struct {
	Kode           string `gorm:"column:Kode;primaryKey;type:varchar(10)" json:"kode"`
	Nama           string `gorm:"column:Nama;type:varchar(100)"          json:"nama"`
	GroupDinasLuar string `gorm:"column:GroupDinasLuar;type:varchar(50)" json:"group_dinas_luar"`
}

// TableName mendefinisikan nama tabel secara eksplisit
func (Jabatan) TableName() string {
	// return "jabatan"
	return "adi_persons.dbo.jabatan"
}
