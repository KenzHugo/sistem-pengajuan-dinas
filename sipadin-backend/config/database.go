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
	// Sesuaikan user, password, host, port, dan dbname milikmu
	dsn := "root:@tcp(127.0.0.1:3306)/sipadin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database: ", err)
	}

	// Auto Migrate Model ke Tabel
	db.AutoMigrate(&models.User{}, &models.PengajuanDinas{}, &models.NotifikasiDireksi{})

	DB = db
	fmt.Println("Koneksi Database Berhasil!")
}