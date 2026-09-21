package dto

import "sipadin/models"

// CreateSppdRequest adalah body untuk membuat pengajuan SPPD baru
type CreateSppdRequest struct {
	KotaTujuan     string           `json:"kota_tujuan"     validate:"required"`
	InstansiTujuan string           `json:"instansi_tujuan" validate:"required"`
	MaksudTujuan   string           `json:"maksud_tujuan"   validate:"required"`
	TglBerangkat   string           `json:"tgl_berangkat"   validate:"required"` // format: YYYY-MM-DD
	TglKembali     string           `json:"tgl_kembali"     validate:"required"` // format: YYYY-MM-DD
	JamBerangkat   string           `json:"jam_berangkat"`                       // format: HH:mm (contoh: 08:00)
	JamKembali     string           `json:"jam_kembali"`                         // format: HH:mm (contoh: 17:00)
	JarakKM        float64          `json:"jarak_km"`                            // estimasi jarak dalam KM
	Transportasi   string           `json:"transportasi"`
	Akomodasi      string           `json:"akomodasi"`
	BiayaMakan     float64          `json:"biaya_makan"`
	BiayaTransport float64          `json:"biaya_transport"`
	BiayaUangSaku  float64          `json:"biaya_uang_saku"`
	BiayaTiket     float64          `json:"biaya_tiket"`
	BiayaHotel     float64          `json:"biaya_hotel"`
	EstimasiBiaya  float64          `json:"estimasi_biaya"`
	RincianBiaya   string           `json:"rincian_biaya"`
	LampiranFile   string           `json:"lampiran_file"` // opsional, path file dokumen pendukung
	Anggota        []AnggotaRequest `json:"anggota"`       // opsional, untuk rombongan
	RekanKerja     []string         `json:"rekan_kerja"`   // opsional, list nomor ID / nama rekan sekantor
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
	JamBerangkat   string           `json:"jam_berangkat"`
	JamKembali     string           `json:"jam_kembali"`
	JarakKM        float64          `json:"jarak_km"`
	Transportasi   string           `json:"transportasi"`
	Akomodasi      string           `json:"akomodasi"`
	BiayaMakan     float64          `json:"biaya_makan"`
	BiayaTransport float64          `json:"biaya_transport"`
	BiayaUangSaku  float64          `json:"biaya_uang_saku"`
	BiayaTiket     float64          `json:"biaya_tiket"`
	BiayaHotel     float64          `json:"biaya_hotel"`
	EstimasiBiaya  float64          `json:"estimasi_biaya"`
	RincianBiaya   string           `json:"rincian_biaya"`
	LampiranFile   string           `json:"lampiran_file"`
	Anggota        []AnggotaRequest `json:"anggota"`
	RekanKerja     []string         `json:"rekan_kerja"`
}

// SppdListFilter adalah filter untuk query daftar SPPD
type SppdListFilter struct {
	Status    string `query:"status"`
	IDPemohon int    `query:"id_pemohon"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
}

// KalkulasiBiayaRequest payload untuk menghitung estimasi biaya dan batas plafon
type KalkulasiBiayaRequest struct {
	NomorID      string   `json:"nomor_id"`
	KotaTujuan   string   `json:"kota_tujuan"`
	TglBerangkat string   `json:"tgl_berangkat"`
	TglKembali   string   `json:"tgl_kembali"`
	JamBerangkat string   `json:"jam_berangkat"`
	JamKembali   string   `json:"jam_kembali"`
	JarakKM      float64  `json:"jarak_km"`
	Transportasi string   `json:"transportasi"`
	Akomodasi    string   `json:"akomodasi"`
	RekanKerja   []string `json:"rekan_kerja"`
}

// RincianPerKepalaDTO rincian hak biaya per personil dalam rombongan
type RincianPerKepalaDTO struct {
	NomorID        string  `json:"nomor_id"`
	Nama           string  `json:"nama"`
	Role           string  `json:"role"` // "Pemohon" | "Rombongan"
	KodeJabatan    string  `json:"kode_jabatan"`
	NamaJabatan    string  `json:"nama_jabatan"`
	GroupDinasLuar string  `json:"group_dinas_luar"`
	Makan          float64 `json:"makan"`
	Saku           float64 `json:"saku"`
	Hotel          float64 `json:"hotel"`
	Transport      float64 `json:"transport"`
	Tiket          float64 `json:"tiket"`
	Subtotal       float64 `json:"subtotal"`
}

// KomponenBersamaDTO rincian biaya bersama (transportasi / tiket rombongan)
type KomponenBersamaDTO struct {
	Transport float64 `json:"transport"`
	Tiket     float64 `json:"tiket"`
	Total     float64 `json:"total"`
}

// KalkulasiBiayaResponse rincian breakdown kalkulasi biaya SPPD dan batas plafon hak
type KalkulasiBiayaResponse struct {
	TotalKeseluruhan    float64               `json:"total_keseluruhan"`
	RincianPerKepala    []RincianPerKepalaDTO `json:"rincian_per_kepala"`
	KomponenBersama     KomponenBersamaDTO    `json:"komponen_bersama"`

	DurasiHari          int     `json:"durasi_hari"`
	DurasiJam           float64 `json:"durasi_jam"`
	JarakKM             float64 `json:"jarak_km"`
	JumlahOrang         int     `json:"jumlah_orang"`
	Area                string  `json:"area"`
	KategoriJabatan     string  `json:"kategori_jabatan"`
	KodeJabatan         string  `json:"kode_jabatan"`
	MataUang            string  `json:"mata_uang"`
	UangMakanPerHari    float64 `json:"uang_makan_per_hari"`
	UangSakuPerHari     float64 `json:"uang_saku_per_hari"`
	UangHarianPerHari   float64 `json:"uang_harian_per_hari"`
	HotelPerMalam       float64 `json:"hotel_per_malam"`
	TransportPerOrang   float64 `json:"transport_per_orang"`
	TarifTransportPerKM float64 `json:"tarif_transport_per_km"`

	// Plafon Maksimal Komponen
	PlafonMakan     float64 `json:"plafon_makan"`
	PlafonTransport float64 `json:"plafon_transport"`
	PlafonUangSaku  float64 `json:"plafon_uang_saku"`
	PlafonTiket     float64 `json:"plafon_tiket"`
	PlafonHotel     float64 `json:"plafon_hotel"`
	PlafonTotal     float64 `json:"plafon_total"`

	// Hak Fasilitas Flags (1 = berhak, 0 = tidak berhak)
	HakTransport int16 `json:"hak_transport"`
	HakUangMakan int16 `json:"hak_uang_makan"`
	HakUangSaku  int16 `json:"hak_uang_saku"`
	HakTiket     int16 `json:"hak_tiket"`
	HakHotel     int16 `json:"hak_hotel"`

	TotalUangHarian    float64 `json:"total_uang_harian"`
	TotalHotel         float64 `json:"total_hotel"`
	TotalTransport     float64 `json:"total_transport"`
	EstimasiBiayaTotal float64 `json:"estimasi_biaya_total"`
	RincianTeks        string  `json:"rincian_teks"`
}

// SppdLedgerResponse payload komprehensif audit ledger keuangan SPPD
type SppdLedgerResponse struct {
	IDSppd             int         `json:"id_sppd"`
	NoSppd             string      `json:"no_sppd"`
	KotaTujuan         string      `json:"kota_tujuan"`
	InstansiTujuan     string      `json:"instansi_tujuan"`
	MaksudTujuan       string      `json:"maksud_tujuan"`
	TglBerangkat       string      `json:"tgl_berangkat"`
	TglKembali         string      `json:"tgl_kembali"`
	DurasiHari         int         `json:"durasi_hari"`
	DurasiJam          float64     `json:"durasi_jam"`
	Transportasi       string      `json:"transportasi"`
	Akomodasi          string      `json:"akomodasi"`
	BiayaTransport     float64     `json:"biaya_transport"`
	BiayaTiket         float64     `json:"biaya_tiket"`
	EstimasiBiayaTotal float64     `json:"estimasi_biaya_total"`
	StatusDokumen      string      `json:"status_dokumen"`
	TotalHakIndividu   float64     `json:"total_hak_individu"`
	TotalBersama       float64     `json:"total_bersama"`
	DetailLedger       []models.SppdDetailBiaya `json:"detail_ledger"`
}

