package dto

// LoginRequest adalah body request untuk endpoint POST /api/auth/login
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // nomor_id atau nik_pendek
	Password   string `json:"password"   validate:"required"`
}

// LoginResponse adalah data yang dikembalikan saat login berhasil
type LoginResponse struct {
	Token       string      `json:"token"`
	User        UserProfile `json:"user"`
}

// UserProfile adalah data publik user yang aman dikirim ke frontend
type UserProfile struct {
	IDUser      int    `json:"id_user"`
	NomorID     string `json:"nomor_id"`
	NamaLengkap string `json:"nama_lengkap"`
	RoleSistem  string `json:"role_sistem"`
	NoWhatsapp  string `json:"no_whatsapp"`
}

// UserClaims adalah payload yang disimpan di dalam JWT token
type UserClaims struct {
	IDUser      int    `json:"id_user"`
	NamaLengkap string `json:"nama_lengkap"`
	RoleSistem  string `json:"role_sistem"`
}
