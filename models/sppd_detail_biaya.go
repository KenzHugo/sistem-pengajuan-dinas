package models

import "time"

// SppdDetailBiaya merepresentasikan tabel `sppd_detail_biaya` sebagai Audit Financial Snapshot Ledger.
// Tabel ini merekam status profil karyawan dan alokasi nominal uang secara permanen (immutable snapshot)
// pada saat persetujuan keuangan agar riwayat audit tidak berubah meskipun master karyawan berubah di masa depan.
type SppdDetailBiaya struct {
	IDDetail            int64     `gorm:"column:id_detail;primaryKey;autoIncrement" json:"id_detail"`
	IDPengajuan         int       `gorm:"column:id_pengajuan;not null;index:idx_ledger_pengajuan" json:"id_pengajuan"`
	IDKaryawan          string    `gorm:"column:id_karyawan;type:nvarchar(50);not null;index:idx_ledger_karyawan" json:"id_karyawan"`
	RolePersonil        string    `gorm:"column:role_personil;type:nvarchar(20);not null;default:'Pemohon';index:idx_ledger_role" json:"role_personil"`
	SnapshotNama        string    `gorm:"column:snapshot_nama;type:nvarchar(128);not null" json:"snapshot_nama"`
	SnapshotJabatan     string    `gorm:"column:snapshot_jabatan;type:nvarchar(100);not null" json:"snapshot_jabatan"`
	SnapshotGrupDinas   string    `gorm:"column:snapshot_grup_dinas;type:nvarchar(50);not null;index:idx_ledger_grup_dinas" json:"snapshot_grup_dinas"`
	SnapshotTujuan      string    `gorm:"column:snapshot_tujuan;type:nvarchar(255);not null" json:"snapshot_tujuan"`
	SnapshotDeskripsi   string    `gorm:"column:snapshot_deskripsi;type:nvarchar(max);not null" json:"snapshot_deskripsi"`
	UangMakan           float64   `gorm:"column:uang_makan;type:decimal(15,2);default:0" json:"uang_makan"`
	UangSaku            float64   `gorm:"column:uang_saku;type:decimal(15,2);default:0" json:"uang_saku"`
	UangHotel           float64   `gorm:"column:uang_hotel;type:decimal(15,2);default:0" json:"uang_hotel"`
	SubtotalHakIndividu float64   `gorm:"column:subtotal_hak_individu;type:decimal(15,2);default:0" json:"subtotal_hak_individu"`
	CatatanKoreksi      string    `gorm:"column:catatan_koreksi;type:nvarchar(max)" json:"catatan_koreksi"`
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime2;autoCreateTime;index:idx_ledger_created_at" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime2;autoUpdateTime" json:"updated_at"`
}

// TableName menentukan nama tabel database secara eksplisit
func (SppdDetailBiaya) TableName() string {
	return "sppd_detail_biaya"
}
