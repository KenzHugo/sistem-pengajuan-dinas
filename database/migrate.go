package database

import (
	"log"
	"sipadin/config"
	"sipadin/models"
)

// Migrate menjalankan auto-migration untuk tabel transaksi sistem
func Migrate() {
	db := config.DB

	err := db.AutoMigrate(
		&models.User{},
		&models.SppdPengajuan{},
		&models.SppdApprovalLog{},
		&models.SppdAnggota{},
		&models.SppdDetailBiaya{},
	)
	if err != nil {
		log.Fatalf("[Migrate] Gagal auto-migrate: %v", err)
	}

	log.Println("[Migrate] Auto-migration selesai untuk tabel master & transaksi")
}
