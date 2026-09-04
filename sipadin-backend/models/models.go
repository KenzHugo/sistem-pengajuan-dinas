package models

import "time"

// ============================================================
// MASTER TABLE: Karyawan_copy1 (pre-existing, read-only via GORM)
// ============================================================

type KaryawanCopy1 struct {
	NomorID       string `gorm:"column:NomorID;primaryKey"         json:"nomor_id"`
	Nama          string `gorm:"column:Nama"                       json:"nama"`
	KodeDepartemen string `gorm:"column:KodeDepartemen"            json:"kode_departemen"`
	KodeJabatan   string `gorm:"column:KodeJabatan"               json:"kode_jabatan"`
	KodeGolongan  string `gorm:"column:KodeGolongan"              json:"kode_golongan"`
	Email         string `gorm:"column:Email"                     json:"email"`
	Password      string `gorm:"column:Password"                  json:"password,omitempty"`
	Aktif         int    `gorm:"column:Aktif"                     json:"aktif"`
}

func (KaryawanCopy1) TableName() string {
	return "Karyawan_copy1"
}

// ============================================================
// MASTER TABLE: PerjalananDinas_Tarif (pre-existing, read-only)
// Composite PK: Area + NoUrut + Jabatan
// ============================================================

type PerjalananDinasTarif struct {
	Area      string  `gorm:"column:Area;primaryKey"      json:"area"`
	NoUrut    float64 `gorm:"column:NoUrut;primaryKey"    json:"no_urut"`
	Jabatan   string  `gorm:"column:Jabatan;primaryKey"   json:"jabatan"`
	Transport float64 `gorm:"column:Transport"            json:"transport"`
	UangMakan float64 `gorm:"column:UangMakan"            json:"uang_makan"`
	UangSaku  float64 `gorm:"column:UangSaku"             json:"uang_saku"`
	Tiket     float64 `gorm:"column:Tiket"                json:"tiket"`
	Hotel     float64 `gorm:"column:Hotel"                json:"hotel"`
	MataUang  string  `gorm:"column:MataUang"             json:"mata_uang"`
}

func (PerjalananDinasTarif) TableName() string {
	return "PerjalananDinas_Tarif"
}

// ============================================================
// MASTER TABLE: PerjalananDinas_Fasilitas (pre-existing, read-only)
// ============================================================

type PerjalananDinasFasilitas struct {
	Kriteria     string `gorm:"column:Kriteria;primaryKey"  json:"kriteria"`
	HakTransport int    `gorm:"column:HakTransport"         json:"hak_transport"`
	HakUangMakan int    `gorm:"column:HakUangMakan"         json:"hak_uang_makan"`
	HakUangSaku  int    `gorm:"column:HakUangSaku"          json:"hak_uang_saku"`
	HakTiket     int    `gorm:"column:HakTiket"             json:"hak_tiket"`
	HakHotel     int    `gorm:"column:HakHotel"             json:"hak_hotel"`
}

func (PerjalananDinasFasilitas) TableName() string {
	return "PerjalananDinas_Fasilitas"
}

// ============================================================
// TRANSACTIONAL TABLE: PerjalananDinas_Pengajuan
// StatusFlow pipeline:
//   SUBMITTED → APPROVED_ATASAN → CONTROLLED_HRD →
//   WA_SENT_DIREKSI → CONFIRMED_DIREKSI → SURAT_TUGAS_ISSUED
//   (any stage → REJECTED)
// ============================================================

type PerjalananDinasPengajuan struct {
	NoPengajuan              string     `gorm:"column:NoPengajuan;primaryKey"          json:"no_pengajuan"`
	NomorID                  string     `gorm:"column:NomorID;not null"                json:"nomor_id"`
	AtasanNomorID            string     `gorm:"column:AtasanNomorID"                   json:"atasan_nomor_id"`
	HRDNomorID               string     `gorm:"column:HRDNomorID"                      json:"hrd_nomor_id"`
	AreaTujuan               string     `gorm:"column:AreaTujuan"                      json:"area_tujuan"`
	KotaTujuan               string     `gorm:"column:KotaTujuan"                      json:"kota_tujuan"`
	MaksudTujuan             string     `gorm:"column:MaksudTujuan;type:text"          json:"maksud_tujuan"`
	KriteriaFasilitas        string     `gorm:"column:KriteriaFasilitas"               json:"kriteria_fasilitas"`
	TanggalBerangkat         time.Time  `gorm:"column:TanggalBerangkat"                json:"tanggal_berangkat"`
	TanggalKembali           time.Time  `gorm:"column:TanggalKembali"                  json:"tanggal_kembali"`
	JamBerangkat             string     `gorm:"column:JamBerangkat"                    json:"jam_berangkat"`
	JamKembali               string     `gorm:"column:JamKembali"                      json:"jam_kembali"`
	EstimasiBiaya            float64    `gorm:"column:EstimasiBiaya;type:decimal(15,2)" json:"estimasi_biaya"`
	StatusFlow               string     `gorm:"column:StatusFlow;type:varchar(50);default:'SUBMITTED'" json:"status_flow"`
	CatatanHRD               string     `gorm:"column:CatatanHRD;type:text"            json:"catatan_hrd"`
	KonfirmasiDireksiNote    string     `gorm:"column:KonfirmasiDireksiNote;type:text" json:"konfirmasi_direksi_note"`
	TanggalKonfirmasiDireksi *time.Time `gorm:"column:TanggalKonfirmasiDireksi"        json:"tanggal_konfirmasi_direksi"`
	NomorSuratTugas          string     `gorm:"column:NomorSuratTugas"                 json:"nomor_surat_tugas"`
	CreatedAt                time.Time  `gorm:"column:CreatedAt;autoCreateTime"        json:"created_at"`
	UpdatedAt                time.Time  `gorm:"column:UpdatedAt;autoUpdateTime"        json:"updated_at"`

	// Virtual associations and enriched fields (not stored in DB columns)
	Karyawan       *KaryawanCopy1 `gorm:"-" json:"karyawan,omitempty"`
	Atasan         *KaryawanCopy1 `gorm:"-" json:"atasan,omitempty"`
	Nama           string         `gorm:"-" json:"nama,omitempty"`
	KodeDepartemen string         `gorm:"-" json:"kode_departemen,omitempty"`
	KodeJabatan    string         `gorm:"-" json:"kode_jabatan,omitempty"`
	AtasanNama     string         `gorm:"-" json:"atasan_nama,omitempty"`
}

func (PerjalananDinasPengajuan) TableName() string {
	return "PerjalananDinas_Pengajuan"
}