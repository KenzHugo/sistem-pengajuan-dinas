package dto

// CreateUserRequest untuk membuat akun user baru
type CreateUserRequest struct {
	NomorID        string `json:"nomor_id"        validate:"required"`
	NikPendek      string `json:"nik_pendek"`
	NamaLengkap    string `json:"nama_lengkap"    validate:"required"`
	KodeDepartemen string `json:"kode_departemen"`
	KodeJabatan    string `json:"kode_jabatan"`
	RoleSistem     string `json:"role_sistem"     validate:"required"`
	NoWhatsapp     string `json:"no_whatsapp"`
	Password       string `json:"password"        validate:"required"`
	StatusAktif    int8   `json:"status_aktif"`
}

// UpdateUserRequest untuk update data user (semua field opsional)
type UpdateUserRequest struct {
	NikPendek      string `json:"nik_pendek"`
	NamaLengkap    string `json:"nama_lengkap"`
	KodeDepartemen string `json:"kode_departemen"`
	KodeJabatan    string `json:"kode_jabatan"`
	RoleSistem     string `json:"role_sistem"`
	NoWhatsapp     string `json:"no_whatsapp"`
	Password       string `json:"password"` // opsional, jika kosong tidak diubah
	StatusAktif    *int8  `json:"status_aktif"`
}
