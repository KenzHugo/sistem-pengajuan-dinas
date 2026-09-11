package dto

// CreateTarifRequest untuk menambah data master tarif fasilitas
type CreateTarifRequest struct {
	AreaTujuan      string  `json:"area_tujuan"      validate:"required"`
	KriteriaJabatan string  `json:"kriteria_jabatan" validate:"required"`
	HakTransport    float64 `json:"hak_transport"`
	HakHotel        float64 `json:"hak_hotel"`
	TarifUangMakan  float64 `json:"tarif_uang_makan"`
	TarifUangSaku   float64 `json:"tarif_uang_saku"`
}

// UpdateTarifRequest untuk update data tarif (semua field opsional)
type UpdateTarifRequest struct {
	AreaTujuan      string  `json:"area_tujuan"`
	KriteriaJabatan string  `json:"kriteria_jabatan"`
	HakTransport    float64 `json:"hak_transport"`
	HakHotel        float64 `json:"hak_hotel"`
	TarifUangMakan  float64 `json:"tarif_uang_makan"`
	TarifUangSaku   float64 `json:"tarif_uang_saku"`
}
