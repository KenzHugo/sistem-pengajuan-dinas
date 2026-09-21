package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config menyimpan seluruh konfigurasi aplikasi yang dibaca dari .env
type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTExpireHours int
	AppPort        string
	WAGatewayURL   string
	WAToken        string
	WAApiKey       string
	WAProvider     string
}

// AppConfig adalah instance global konfigurasi
var AppConfig *Config

// LoadConfig membaca file .env dan mengisi AppConfig
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		if err2 := godotenv.Load(".env"); err2 != nil {
			log.Println("[Config] File .env tidak ditemukan, membaca default & OS env")
		}
	}

	expireHours, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	if err != nil {
		expireHours = 24
	}

	AppConfig = &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "sipadin"),
		JWTSecret:      getEnv("JWT_SECRET", "secret-fallback"),
		JWTExpireHours: expireHours,
		AppPort:        getEnv("APP_PORT", "3000"),
		WAGatewayURL:   getEnv("WA_GATEWAY_URL", ""),
		WAToken:        getEnv("WA_TOKEN", getEnv("FONNTE_TOKEN", "")),
		WAApiKey:       getEnv("WA_API_KEY", ""),
		WAProvider:     getEnv("WA_PROVIDER", "fonnte"),
	}

	log.Println("[Config] Konfigurasi berhasil dimuat")
}

// getEnv mengambil nilai env dengan fallback jika tidak ada
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
