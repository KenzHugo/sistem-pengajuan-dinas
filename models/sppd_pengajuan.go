package models

import "time"

// SppdPengajuan merepresentasikan tabel `sppd_pengajuan` — master dokumen SPPD.
// Kolom persis sesuai PRD v2.0. TIDAK ada kolom tambahan.
type SppdPengajuan struct {
	IDSppd         int       `gorm:"column:id_sppd;primaryKey;autoIncrement"              json:"id_sppd"`
	NoSppd         string    `gorm:"column:no_sppd;type:varchar(50)"                      json:"no_sppd"`
	IDPemohon      int       `gorm:"column:id_pemohon;type:int;not null"                  json:"id_pemohon"`
	KotaTujuan     string    `gorm:"column:kota_tujuan;type:varchar(100);not null"        json:"kota_tujuan"`
	InstansiTujuan string    `gorm:"column:instansi_tujuan;type:varchar(150);not null"    json:"instansi_tujuan"`
	MaksudTujuan   string    `gorm:"column:maksud_tujuan;type:varchar(max);not null"      json:"maksud_tujuan"`
	TglBerangkat   time.Time `gorm:"column:tgl_berangkat;type:date;not null"              json:"tgl_berangkat"`
	TglKembali     time.Time `gorm:"column:tgl_kembali;type:date;not null"                json:"tgl_kembali"`
	JamBerangkat   string    `gorm:"column:jam_berangkat;type:varchar(10);default:'08:00'" json:"jam_berangkat"`
	JamKembali     string    `gorm:"column:jam_kembali;type:varchar(10);default:'17:00'"   json:"jam_kembali"`
	DurasiHari     int       `gorm:"column:durasi_hari"                                   json:"durasi_hari"`
	DurasiJam      float64   `gorm:"column:durasi_jam;type:decimal(10,2);default:0"       json:"durasi_jam"`
	JarakKM        float64   `gorm:"column:jarak_km;type:decimal(10,2);default:0"         json:"jarak_km"`
	Transportasi   string    `gorm:"column:transportasi;type:varchar(50)"                 json:"transportasi"`
	Akomodasi      string    `gorm:"column:akomodasi;type:varchar(50)"                    json:"akomodasi"`
	BiayaMakan     float64   `gorm:"column:biaya_makan;type:decimal(20,4);default:0"      json:"biaya_makan"`
	BiayaTransport float64   `gorm:"column:biaya_transport;type:decimal(20,4);default:0"  json:"biaya_transport"`
	BiayaUangSaku  float64   `gorm:"column:biaya_uang_saku;type:decimal(20,4);default:0"   json:"biaya_uang_saku"`
	BiayaTiket     float64   `gorm:"column:biaya_tiket;type:decimal(20,4);default:0"      json:"biaya_tiket"`
	BiayaHotel     float64   `gorm:"column:biaya_hotel;type:decimal(20,4);default:0"      json:"biaya_hotel"`
	PlafonMakan     float64   `gorm:"column:plafon_makan;type:decimal(20,4);default:0"     json:"plafon_makan"`
	PlafonTransport float64   `gorm:"column:plafon_transport;type:decimal(20,4);default:0" json:"plafon_transport"`
	PlafonUangSaku  float64   `gorm:"column:plafon_uang_saku;type:decimal(20,4);default:0"  json:"plafon_uang_saku"`
	PlafonTiket     float64   `gorm:"column:plafon_tiket;type:decimal(20,4);default:0"     json:"plafon_tiket"`
	PlafonHotel     float64   `gorm:"column:plafon_hotel;type:decimal(20,4);default:0"     json:"plafon_hotel"`
	StatusPlafon    string    `gorm:"column:status_plafon;type:varchar(50);default:'SESUAI'" json:"status_plafon"`
	EstimasiBiaya  float64   `gorm:"column:estimasi_biaya;type:decimal(20,4)"             json:"estimasi_biaya"`
	RincianBiaya   string    `gorm:"column:rincian_biaya;type:varchar(max)"               json:"rincian_biaya"`
	LampiranFile   string    `gorm:"column:lampiran_file;type:varchar(255)"               json:"lampiran_file"`
	StatusDokumen  string    `gorm:"column:status_dokumen;type:varchar(50);default:draft" json:"status_dokumen"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"                     json:"created_at"`

	// Relasi (preload only, tidak menjadi kolom di DB)
	Pemohon         *User             `gorm:"foreignKey:IDPemohon;references:IDUser"  json:"pemohon,omitempty"`
	ApprovalLogs    []SppdApprovalLog `gorm:"foreignKey:IDSppd;references:IDSppd"     json:"approval_logs,omitempty"`
	Anggota         []SppdAnggota     `gorm:"foreignKey:IDSppd;references:IDSppd"     json:"anggota,omitempty"`
	DetailBiaya     []SppdDetailBiaya `gorm:"foreignKey:IDPengajuan;references:IDSppd" json:"detail_biaya,omitempty"`
	NamaManagerHRGA string            `gorm:"-"                                       json:"nama_manager_hrga,omitempty"`
	NamaDireksi     string            `gorm:"-"                                       json:"nama_direksi,omitempty"`
}

func (SppdPengajuan) TableName() string {
	return "sppd_pengajuan"
}
