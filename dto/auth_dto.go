package dto

// LoginRequest adalah body request untuk endpoint POST /api/auth/login
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // nomor_id atau nik_pendek
	Password   string `json:"password"   validate:"required"`
	ClientTime string `json:"client_time"` // Waktu ISO/RFC3339 dari device pengguna
}

// LoginResponse adalah data yang dikembalikan saat login berhasil
type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

// UserProfile adalah data publik user yang aman dikirim ke frontend
type UserProfile struct {
	IDUser          int    `json:"id_user"`
	NomorID         string `json:"nomor_id"`
	NikPendek       string `json:"nik_pendek"`
	NamaLengkap     string `json:"nama_lengkap"`
	RoleSistem      string `json:"role_sistem"`
	NoWhatsapp      string `json:"no_whatsapp"`
	IsSpecialAccess bool   `json:"is_special_access"`
}

// UserClaims adalah payload yang disimpan di dalam JWT token
type UserClaims struct {
	IDUser          int    `json:"id_user"`
	NomorID         string `json:"nomor_id"`
	NikPendek       string `json:"nik_pendek"`
	NamaLengkap     string `json:"nama_lengkap"`
	RoleSistem      string `json:"role_sistem"`
	IsSpecialAccess bool   `json:"is_special_access"`
}

// OperationalStatusResponse adalah response status jam operasional aplikasi dengan cross-check server & device
type OperationalStatusResponse struct {
	IsOpen      bool   `json:"is_open"`
	Message     string `json:"message"`
	Hari        string `json:"hari"`
	JamMasuk    string `json:"jam_masuk"`
	JamSelesai  string `json:"jam_selesai"`
	StatusLibur bool   `json:"status_libur"`
	ServerTime  string `json:"server_time"` // Jam server (WIB)
	DeviceTime  string `json:"device_time"` // Jam device (WIB)
	CurrentTime string `json:"current_time"` // Jam server WIB untuk kompatibilitas
	IsSunday    bool   `json:"is_sunday"`
}
