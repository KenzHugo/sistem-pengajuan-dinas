package models

// User merepresentasikan tabel `users` — master data karyawan dan akun login.
// Kolom persis sesuai PRD v2.0. TIDAK ada kolom tambahan.
type User struct {
	IDUser         int    `gorm:"column:id_user;primaryKey;autoIncrement"       json:"id_user"`
	NomorID        string `gorm:"column:nomor_id;type:varchar(50);not null"     json:"nomor_id"`
	NikPendek      string `gorm:"column:nik_pendek;type:varchar(50)"            json:"nik_pendek"`
	NamaLengkap    string `gorm:"column:nama_lengkap;type:varchar(128);not null" json:"nama_lengkap"`
	KodeDepartemen string `gorm:"column:kode_departemen;type:varchar(16)"       json:"kode_departemen"`
	KodeJabatan    string `gorm:"column:kode_jabatan;type:varchar(16)"          json:"kode_jabatan"`
	RoleSistem     string `gorm:"column:role_sistem;type:varchar(20);not null"  json:"role_sistem"`
	NoWhatsapp     string `gorm:"column:no_whatsapp;type:varchar(20)"           json:"no_whatsapp"`
	PasswordHash   string `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	StatusAktif    int8   `gorm:"column:status_aktif;type:tinyint(1);default:1" json:"status_aktif"`
}

// TableName mendefinisikan nama tabel secara eksplisit agar GORM tidak pluralize
func (User) TableName() string {
	return "users"
}
