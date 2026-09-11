package models

// MasterTarifFasilitas merepresentasikan tabel `master_tarif_fasilitas` — referensi tarif.
// Kolom persis sesuai PRD v2.0. TIDAK ada kolom tambahan.
type MasterTarifFasilitas struct {
	IDTarif         int     `gorm:"column:id_tarif;primaryKey;autoIncrement"          json:"id_tarif"`
	AreaTujuan      string  `gorm:"column:area_tujuan;type:varchar(50);not null"      json:"area_tujuan"`
	KriteriaJabatan string  `gorm:"column:kriteria_jabatan;type:varchar(50);not null" json:"kriteria_jabatan"`
	HakTransport    float64 `gorm:"column:hak_transport;type:decimal(20,4)"          json:"hak_transport"`
	HakHotel        float64 `gorm:"column:hak_hotel;type:decimal(20,4)"              json:"hak_hotel"`
	TarifUangMakan  float64 `gorm:"column:tarif_uang_makan;type:decimal(20,4)"       json:"tarif_uang_makan"`
	TarifUangSaku   float64 `gorm:"column:tarif_uang_saku;type:decimal(20,4)"        json:"tarif_uang_saku"`
}

func (MasterTarifFasilitas) TableName() string {
	return "master_tarif_fasilitas"
}
