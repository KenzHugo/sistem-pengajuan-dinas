package dto

// CreateKaryawanRequest merepresentasikan payload untuk membuat data karyawan baru
type CreateKaryawanRequest struct {
	NomorID        string `json:"nomor_id"`
	Nama           string `json:"nama"`
	KodeDepartemen string `json:"kode_departemen"`
	KodeJabatan    string `json:"kode_jabatan"`
	Job            string `json:"job"`
	Kota           string `json:"kota"`
	TanggalLahir   string `json:"tanggal_lahir"` // Format: YYYY-MM-DD atau kosong
	Aktif          *int   `json:"aktif"`
}

// UpdateKaryawanRequest merepresentasikan payload untuk memperbarui data karyawan
type UpdateKaryawanRequest struct {
	Nama           string `json:"nama"`
	KodeDepartemen string `json:"kode_departemen"`
	KodeJabatan    string `json:"kode_jabatan"`
	Job            string `json:"job"`
	Kota           string `json:"kota"`
	TanggalLahir   string `json:"tanggal_lahir"` // Format: YYYY-MM-DD atau kosong
	Aktif          *int   `json:"aktif"`
}

// KaryawanQuery merepresentasikan parameter filter dan pagination untuk daftar karyawan
type KaryawanQuery struct {
	Search         string `query:"search"`
	KodeDepartemen string `query:"kode_departemen"`
	Aktif          string `query:"aktif"` // "1", "0", atau "" (semua)
	Page           int    `query:"page"`
	Limit          int    `query:"limit"`
	SortBy         string `query:"sort_by"`
	Order          string `query:"order"`
}

// KaryawanResponse merepresentasikan struktur respons data karyawan
type KaryawanResponse struct {
	NomorID              string  `json:"nomor_id"`
	Nama                 string  `json:"nama"`
	KodeDepartemen       string  `json:"kode_departemen"`
	KodeJabatan          string  `json:"kode_jabatan"`
	Job                  string  `json:"job"`
	Kota                 string  `json:"kota"`
	TanggalLahir         *string `json:"tanggal_lahir"`
	TanggalLahirFormatted string  `json:"tanggal_lahir_formatted"`
	Aktif                int     `json:"aktif"`
}

// KaryawanListResponse merepresentasikan respons daftar karyawan beserta pagination & metadata
type KaryawanListResponse struct {
	Total       int64              `json:"total"`
	Page        int                `json:"page"`
	Limit       int                `json:"limit"`
	TotalPages  int                `json:"total_pages"`
	Data        []KaryawanResponse `json:"data"`
	Departemens []string           `json:"departemens,omitempty"`
}

// KaryawanStatsResponse merepresentasikan ringkasan statistik karyawan
type KaryawanStatsResponse struct {
	TotalKaryawan   int64 `json:"total_karyawan"`
	KaryawanAktif   int64 `json:"karyawan_aktif"`
	KaryawanNonAktif int64 `json:"karyawan_non_aktif"`
	TotalDepartemen int64 `json:"total_departemen"`
}
