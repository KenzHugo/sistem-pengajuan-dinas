package main

import (
	"fmt"
	"log"

	"sipadin/config"
	"sipadin/controller"
	"sipadin/database"
	"sipadin/repository"
	"sipadin/routes"
	"sipadin/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// ── 1. Load Konfigurasi ──────────────────────────────────────────────────
	config.LoadConfig()

	// ── 2. Inisialisasi Database ─────────────────────────────────────────────
	config.InitDatabase()

	// ── 3. Auto-migrate (Skema DB) ───────────────────────────────────────────
	database.Migrate()

	// ── 4. Inisialisasi Fiber ────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName:      "SIPADIN - Sistem Pengajuan Perjalanan Dinas v2.0",
		ErrorHandler: globalErrorHandler,
	})

	// ── 5. Global Middleware ─────────────────────────────────────────────────
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} ${method} ${path} — ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS, PATCH",
	}))

	// ── 6. Dependency Injection (Wire semua layer) ────────────────────────────
	// Repositories
	userRepo     := repository.NewUserRepository(config.DB)
	sppdRepo     := repository.NewSppdRepository(config.DB)
	approvalRepo := repository.NewApprovalRepository(config.DB)
	anggotaRepo  := repository.NewAnggotaRepository(config.DB)
	tarifRepo    := repository.NewTarifRepository(config.DB)
	karyawanRepo := repository.NewKaryawanRepository(config.DB)

	// Services
	notifSvc    := service.NewNotificationService(sppdRepo, userRepo)
	authSvc     := service.NewAuthService(userRepo)
	sppdSvc     := service.NewSppdService(sppdRepo, anggotaRepo, tarifRepo, userRepo, notifSvc)
	approvalSvc := service.NewApprovalService(sppdRepo, approvalRepo, userRepo, notifSvc, sppdSvc)
	tarifSvc    := service.NewTarifService(tarifRepo)
	userSvc     := service.NewUserService(userRepo)
	karyawanSvc := service.NewKaryawanService(karyawanRepo)

	// Controllers
	authCtrl     := controller.NewAuthController(authSvc)
	sppdCtrl     := controller.NewSppdController(sppdSvc)
	approvalCtrl := controller.NewApprovalController(approvalSvc)
	tarifCtrl    := controller.NewTarifController(tarifSvc)
	userCtrl     := controller.NewUserController(userSvc)
	karyawanCtrl := controller.NewKaryawanController(karyawanSvc)

	// ── 7. Register API Routes (HARUS sebelum Static) ─────────────────────────
	// PENTING: routes didaftarkan lebih dulu agar POST /api/* tidak dicegat
	// oleh static middleware yang hanya mengizinkan method GET.
	routes.RegisterRoutes(app, authCtrl, sppdCtrl, approvalCtrl, tarifCtrl, userCtrl, karyawanCtrl)

	// ── 8. Static Files (Frontend) — didaftarkan TERAKHIR ──────────────────────
	app.Static("/", "./public", fiber.Static{
		Index: "index.html",
	})

	// ── 9. Jalankan Server ────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", config.AppConfig.AppPort)
	log.Printf("[SIPADIN] Server berjalan di http://localhost%s", addr)
	log.Fatal(app.Listen(addr))
}

// globalErrorHandler menangani panic dan error tidak tertangkap
func globalErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return ctx.Status(code).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
	})
}
