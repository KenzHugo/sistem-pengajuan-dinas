package models

// PerjalananDinasTarif merepresentasikan tabel master tarif `perjalanandinas_tarif`
type PerjalananDinasTarif struct {
	NoUrut    string  `gorm:"column:NoUrut;type:decimal(10,4)" json:"no_urut"`
	Area      string  `gorm:"column:Area;primaryKey;type:varchar(32)" json:"area"`
	Jabatan   string  `gorm:"column:Jabatan;primaryKey;type:varchar(32)" json:"jabatan"`
	MataUang  string  `gorm:"column:MataUang;type:varchar(10)" json:"mata_uang"`
	UangMakan float64 `gorm:"column:UangMakan;type:decimal(20,4)" json:"uang_makan"`
	UangSaku  float64 `gorm:"column:UangSaku;type:decimal(20,4)" json:"uang_saku"`
	Hotel     float64 `gorm:"column:Hotel;type:decimal(20,4)" json:"hotel"`
	Transport float64 `gorm:"column:Transport;type:decimal(20,4)" json:"transport"`
	Tiket     float64 `gorm:"column:Tiket;type:decimal(20,4)" json:"tiket"`
}

func (PerjalananDinasTarif) TableName() string {
	return "perjalanandinas_tarif"
}

// PerjalananDinasFasilitas merepresentasikan tabel ketentuan fasilitas `perjalanandinas_fasilitas`
type PerjalananDinasFasilitas struct {
	Kriteria     string `gorm:"column:Kriteria;primaryKey;type:varchar(128)" json:"kriteria"`
	HakTransport int16  `gorm:"column:HakTransport;type:smallint" json:"hak_transport"`
	HakUangMakan int16  `gorm:"column:HakUangMakan;type:smallint" json:"hak_uang_makan"`
	HakUangSaku  int16  `gorm:"column:HakUangSaku;type:smallint" json:"hak_uang_saku"`
	HakTiket     int16  `gorm:"column:HakTiket;type:smallint" json:"hak_tiket"`
	HakHotel     int16  `gorm:"column:HakHotel;type:smallint" json:"hak_hotel"`
}

func (PerjalananDinasFasilitas) TableName() string {
	return "perjalanandinas_fasilitas"
}

// Alias tipe untuk kompatibilitas
type MasterTarifFasilitas = PerjalananDinasTarif
