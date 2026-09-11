package database

import (
	"log"
	"sipadin/config"
	"sipadin/models"
)

// Migrate menjalankan auto-migration untuk semua 5 tabel sesuai PRD v2.0
func Migrate() {
	db := config.DB

	err := db.AutoMigrate(
		&models.User{},
		&models.SppdPengajuan{},
		&models.SppdApprovalLog{},
		&models.SppdAnggota{},
		&models.MasterTarifFasilitas{},
	)
	if err != nil {
		log.Fatalf("[Migrate] Gagal auto-migrate: %v", err)
	}

	log.Println("[Migrate] Auto-migration selesai untuk 5 tabel")
}
