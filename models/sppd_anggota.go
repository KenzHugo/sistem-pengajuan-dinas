package models

// SppdAnggota merepresentasikan tabel `sppd_anggota` — personil dinas (rombongan).
// Kolom persis sesuai PRD v2.0. TIDAK ada kolom tambahan.
type SppdAnggota struct {
	IDAnggota         int     `gorm:"column:id_anggota;primaryKey;autoIncrement"          json:"id_anggota"`
	IDSppd            int     `gorm:"column:id_sppd;type:int;not null"                             json:"id_sppd"`
	IDUser            int     `gorm:"column:id_user;type:int;not null"                             json:"id_user"`
	KriteriaFasilitas string  `gorm:"column:kriteria_fasilitas;type:varchar(50)"          json:"kriteria_fasilitas"`
	UangMakan         float64 `gorm:"column:uang_makan;type:decimal(20,4);default:0"      json:"uang_makan"`
	UangSaku          float64 `gorm:"column:uang_saku;type:decimal(20,4);default:0"       json:"uang_saku"`

	// Relasi
	User *User `gorm:"foreignKey:IDUser;references:IDUser" json:"user,omitempty"`
}

func (SppdAnggota) TableName() string {
	return "sppd_anggota"
}
