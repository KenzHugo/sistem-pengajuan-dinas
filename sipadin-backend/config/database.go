package config

import (
	"fmt"
	"log"
	"sipadin-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Sesuaikan user, password, host, port, dan dbname
	dsn := "root:@tcp(127.0.0.1:3306)/sipadin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Gagal koneksi ke database: ", err)
	}

	// Drop and AutoMigrate untuk memastikan seluruh kolom dan tipe data sesuai
	_ = db.Migrator().DropTable(&models.NotifikasiDireksi{}, &models.PengajuanDinas{}, &models.User{})
	if err := db.AutoMigrate(&models.User{}, &models.PengajuanDinas{}, &models.NotifikasiDireksi{}); err != nil {
		log.Println("Migrate warning:", err)
	}

	DB = db
	fmt.Println("Koneksi Database Berhasil!")

	// Seed data awal jika tabel kosong
	SeedData()
}