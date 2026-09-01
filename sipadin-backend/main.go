package main

import (
	"log"
	"sipadin-backend/config"
	"sipadin-backend/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// 1. Koneksi Database
	config.ConnectDB()

	// 2. Inisialisasi Fiber
	app := fiber.New()

	// 3. Middleware CORS (Penting agar Frontend HTML/JS lokal bisa akses)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// 4. Setup Routing API
	api := app.Group("/api")

	// Endpoint Auth
	api.Post("/login", controllers.Login)

	// Endpoint Pengajuan Dinas
	api.Get("/pengajuan", controllers.GetPengajuan)
	api.Post("/pengajuan", controllers.CreatePengajuan)
	api.Put("/pengajuan/:id/status", controllers.UpdateStatus)

	// 5. Jalankan Server di Port 8080
	log.Fatal(app.Listen(":8080"))
}