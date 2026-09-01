package main

import (
	"log"
	"sipadin-backend/config"
	"sipadin-backend/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// 1. Koneksi Database & Seed Data
	config.ConnectDB()

	// 2. Inisialisasi Fiber
	app := fiber.New()

	// 3. Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// 4. Setup Routing API
	api := app.Group("/api")

	// Endpoint Auth
	api.Post("/login", controllers.Login)

	// Endpoint Pengajuan Dinas
	api.Get("/pengajuan", controllers.GetPengajuan)
	api.Get("/pengajuan/:id", controllers.GetPengajuanByID)
	api.Post("/pengajuan", controllers.CreatePengajuan)
	api.Put("/pengajuan/:id/status", controllers.UpdateStatus)

	// 5. Static File Serving (Akses langsung web via http://localhost:8080)
	app.Static("/", "../")

	// 6. Jalankan Server di Port 8080
	log.Println("SIPADIN Server berjalan di http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}