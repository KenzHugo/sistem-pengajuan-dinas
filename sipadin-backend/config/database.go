package config

import (
	"fmt"
	"log"
	"sipadin-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	// DSN for sipadin MySQL database
	dsn := "root:@tcp(127.0.0.1:3306)/sipadin?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Disable FK constraints during migration so master tables aren't affected
		DisableForeignKeyConstraintWhenMigrating: true,
		// Suppress verbose SQL logs in production; change to logger.Info for debugging
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal("Gagal koneksi ke database: ", err)
	}

	// Only auto-migrate the transactional table.
	// Master tables (Karyawan_copy1, PerjalananDinas_Tarif, PerjalananDinas_Fasilitas)
	// already exist and are populated from the provided SQL files — do NOT touch them.
	if err := db.AutoMigrate(&models.PerjalananDinasPengajuan{}); err != nil {
		log.Println("AutoMigrate warning:", err)
	}

	DB = db
	fmt.Println("✅ Koneksi Database Berhasil! Tabel PerjalananDinas_Pengajuan siap.")
}