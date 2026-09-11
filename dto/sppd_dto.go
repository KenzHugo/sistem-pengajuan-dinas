package dto

// CreateSppdRequest adalah body untuk membuat pengajuan SPPD baru
type CreateSppdRequest struct {
	KotaTujuan     string           `json:"kota_tujuan"     validate:"required"`
	InstansiTujuan string           `json:"instansi_tujuan" validate:"required"`
	MaksudTujuan   string           `json:"maksud_tujuan"   validate:"required"`
	TglBerangkat   string           `json:"tgl_berangkat"   validate:"required"` // format: YYYY-MM-DD
	TglKembali     string           `json:"tgl_kembali"     validate:"required"` // format: YYYY-MM-DD
	Transportasi   string           `json:"transportasi"`
	Akomodasi      string           `json:"akomodasi"`
	EstimasiBiaya  float64          `json:"estimasi_biaya"`
	Anggota        []AnggotaRequest `json:"anggota"` // opsional, untuk rombongan
}

// AnggotaRequest adalah data satu anggota dalam rombongan dinas
type AnggotaRequest struct {
	IDUser            int     `json:"id_user"            validate:"required"`
	KriteriaFasilitas string  `json:"kriteria_fasilitas"`
	UangMakan         float64 `json:"uang_makan"`
	UangSaku          float64 `json:"uang_saku"`
}

// UpdateSppdRequest untuk mengedit pengajuan yang masih berstatus draft
type UpdateSppdRequest struct {
	KotaTujuan     string           `json:"kota_tujuan"`
	InstansiTujuan string           `json:"instansi_tujuan"`
	MaksudTujuan   string           `json:"maksud_tujuan"`
	TglBerangkat   string           `json:"tgl_berangkat"`
	TglKembali     string           `json:"tgl_kembali"`
	Transportasi   string           `json:"transportasi"`
	Akomodasi      string           `json:"akomodasi"`
	EstimasiBiaya  float64          `json:"estimasi_biaya"`
	Anggota        []AnggotaRequest `json:"anggota"`
}

// SppdListFilter adalah filter untuk query daftar SPPD
type SppdListFilter struct {
	Status    string `query:"status"`
	IDPemohon int    `query:"id_pemohon"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
}
