package config

import (
	"fmt"
	"log"
	"net/url"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB adalah instance global koneksi database GORM
var DB *gorm.DB

// InitDatabase menginisialisasi koneksi ke SQL Server menggunakan GORM
func InitDatabase() {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable&trustServerCertificate=true",
		url.QueryEscape(AppConfig.DBUser),
		url.QueryEscape(AppConfig.DBPassword),
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Warn),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("[Database] Gagal koneksi ke SQL Server: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[Database] Gagal mendapatkan *sql.DB: %v", err)
	}

	// Connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	log.Println("[Database] Koneksi berhasil ke SQL Server database:", AppConfig.DBName)
}
