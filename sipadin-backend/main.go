package main

import (
	"log"
	"sipadin-backend/config"
	"sipadin-backend/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// 1. Connect to MySQL & run AutoMigrate on PerjalananDinas_Pengajuan
	config.ConnectDB()

	// 2. Initialize Fiber
	app := fiber.New(fiber.Config{
		AppName: "SIPADIN Backend v3.0",
	})

	// 3. CORS middleware — allow all origins for local HTML/LiveServer dev
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// 4. API routing
	api := app.Group("/api")

	// ── Auth ──────────────────────────────────────────────────
	auth := api.Group("/auth")
	auth.Post("/login", controllers.Login)

	// ── Master Lookup Data ────────────────────────────────────
	master := api.Group("/master")
	master.Get("/tarif", controllers.GetMasterTarif)
	master.Get("/fasilitas", controllers.GetMasterFasilitas)

	// ── Pengajuan CRUD ────────────────────────────────────────
	pengajuan := api.Group("/pengajuan")
	pengajuan.Post("/", controllers.CreatePengajuan)
	pengajuan.Get("/", controllers.GetPengajuan)
	pengajuan.Get("/detail", controllers.GetPengajuanByID)
	pengajuan.Get("/detail/*", controllers.GetPengajuanByID)

	// Stage 2: Atasan approves or rejects
	pengajuan.Put("/action/atasan-approve", controllers.AtasanApprove)
	pengajuan.Put("/action/atasan-reject", controllers.AtasanReject)
	pengajuan.Put("/action/approve", controllers.AtasanApprove)
	pengajuan.Put("/action/reject", controllers.AtasanReject)
	pengajuan.Put("/:id/atasan-approve", controllers.AtasanApprove)
	pengajuan.Put("/:id/atasan-reject", controllers.AtasanReject)
	pengajuan.Put("/:id/approve", controllers.AtasanApprove)
	pengajuan.Put("/:id/reject", controllers.AtasanReject)

	// Stage 3: HRGA controls & verifies / rejects
	pengajuan.Put("/action/hrd-control", controllers.HRDControl)
	pengajuan.Put("/:id/hrd-control", controllers.HRDControl)
	pengajuan.Put("/action/hrd-reject", controllers.HRDReject)
	pengajuan.Put("/:id/hrd-reject", controllers.HRDReject)

	// Stage 4: HRGA generates WA message → opens WhatsApp
	pengajuan.Get("/action/wa-text", controllers.GetWAText)
	pengajuan.Get("/:id/wa-text", controllers.GetWAText)

	// Stage 5: HRGA records Direksi confirmation / rejection
	pengajuan.Put("/action/direksi-confirm", controllers.DireksiConfirm)
	pengajuan.Put("/:id/direksi-confirm", controllers.DireksiConfirm)
	pengajuan.Put("/action/direksi-reject", controllers.DireksiReject)
	pengajuan.Put("/:id/direksi-reject", controllers.DireksiReject)

	// Stage 6: HRGA issues official Surat Tugas
	pengajuan.Post("/action/surat-tugas", controllers.IssueSuratTugas)
	pengajuan.Put("/action/surat-tugas", controllers.IssueSuratTugas)
	pengajuan.Post("/:id/surat-tugas", controllers.IssueSuratTugas)
	pengajuan.Put("/:id/issue-surat-tugas", controllers.IssueSuratTugas)

	// Fallback ID route
	pengajuan.Get("/:id", controllers.GetPengajuanByID)

	// 5. Static file serving — serve frontend HTML directly via http://localhost:8080
	app.Static("/", "../")

	// 6. Start server
	log.Println("🚀 SIPADIN Server v3.0 berjalan di http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}