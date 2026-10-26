package routes

import (
	"sipadin/controller"
	"sipadin/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes mendaftarkan semua route API ke aplikasi Fiber
func RegisterRoutes(
	app *fiber.App,
	authCtrl *controller.AuthController,
	sppdCtrl *controller.SppdController,
	approvalCtrl *controller.ApprovalController,
	tarifCtrl *controller.TarifController,
	userCtrl *controller.UserController,
	karyawanCtrl *controller.KaryawanController,
) {
	api := app.Group("/api")

	// ─── Auth Routes (public) ──────────────────────────────────────────────────
	auth := api.Group("/auth")
	auth.Post("/login", authCtrl.Login)
	auth.Get("/users", authCtrl.GetDBUsers)

	// ─── Protected Routes (JWT Required) ──────────────────────────────────────
	protected := api.Group("", middleware.Protected())

	// ─── SPPD Routes ──────────────────────────────────────────────────────────
	sppd := protected.Group("/sppd")

	// Semua role yang login bisa akses
	sppd.Get("/mine", sppdCtrl.GetMyList)
	sppd.Get("/queue", sppdCtrl.GetQueue)
	sppd.Get("/:id", sppdCtrl.GetByID)
	sppd.Get("/", sppdCtrl.GetAll)
	sppd.Post("/hitung-biaya", sppdCtrl.HitungBiaya)
	sppd.Post("/upload-lampiran", sppdCtrl.UploadLampiran)

	// Karyawan: membuat, mengedit, menghapus, mengajukan
	sppd.Post("/",
		middleware.RequireRole(middleware.RoleKaryawan, middleware.RoleManager, middleware.RoleHRGA, middleware.RoleFinance),
		sppdCtrl.Create,
	)
	sppd.Post("/:id/submit",
		middleware.RequireRole(middleware.RoleKaryawan, middleware.RoleManager, middleware.RoleHRGA, middleware.RoleFinance),
		sppdCtrl.Submit,
	)
	sppd.Put("/:id",
		middleware.RequireRole(middleware.RoleKaryawan, middleware.RoleManager, middleware.RoleHRGA, middleware.RoleFinance),
		sppdCtrl.Update,
	)
	sppd.Delete("/:id",
		middleware.RequireRole(middleware.RoleKaryawan, middleware.RoleManager, middleware.RoleHRGA, middleware.RoleFinance),
		sppdCtrl.Delete,
	)

	// ─── Approval Routes ───────────────────────────────────────────────────────
	approval := protected.Group("/approval")

	approval.Post("/:id/manager",
		middleware.RequireRole(middleware.RoleManager),
		approvalCtrl.ManagerAction,
	)
	approval.Post("/:id/hrga",
		middleware.RequireRole(middleware.RoleHRGA),
		approvalCtrl.HRGAAction,
	)
	approval.Post("/:id/finance",
		middleware.RequireRole(middleware.RoleFinance),
		approvalCtrl.FinanceAction,
	)
	approval.Post("/:id/send",
		middleware.RequireRole(middleware.RoleHRGA),
		approvalCtrl.SendToDirector,
	)

	// ─── Finance & Audit Ledger Routes ─────────────────────────────────────────
	finance := protected.Group("/finance")
	finance.Get("/sppd-ledger/:id", approvalCtrl.GetSppdLedger)
	finance.Get("/export-csv", approvalCtrl.ExportFinancialCSV)

	// ─── Tarif Routes (HRGA Full Access, semua lainnya GET only) ──────────────
	tarif := protected.Group("/tarif")
	tarif.Get("/fasilitas", tarifCtrl.GetFasilitas)
	tarif.Get("/areas", tarifCtrl.GetAreas)
	tarif.Get("/", tarifCtrl.GetAll)
	tarif.Get("/:id", tarifCtrl.GetByID)
	tarif.Post("/", middleware.RequireRole(middleware.RoleHRGA), tarifCtrl.Create)
	tarif.Put("/", middleware.RequireRole(middleware.RoleHRGA), tarifCtrl.Update)
	tarif.Put("/:id", middleware.RequireRole(middleware.RoleHRGA), tarifCtrl.Update)
	tarif.Delete("/", middleware.RequireRole(middleware.RoleHRGA), tarifCtrl.Delete)
	tarif.Delete("/:id", middleware.RequireRole(middleware.RoleHRGA), tarifCtrl.Delete)

	// ─── User Management Routes (HRGA Full Access) ────────────────────────────
	users := protected.Group("/users")
	users.Get("/", userCtrl.GetAll) // semua bisa lihat (untuk select anggota)
	users.Get("/list", userCtrl.GetAll)
	users.Get("/:id", userCtrl.GetByID)
	users.Post("/", middleware.RequireRole(middleware.RoleHRGA), userCtrl.Create)
	users.Put("/:id", middleware.RequireRole(middleware.RoleHRGA), userCtrl.Update)
	users.Delete("/:id", middleware.RequireRole(middleware.RoleHRGA), userCtrl.Delete)

	// ─── Karyawan Routes (SQL Server adi_persons.dbo.karyawan) ─────────────────
	karyawan := protected.Group("/karyawan", middleware.RequireRole(
		middleware.RoleManager,
		middleware.RoleHRGA,
		middleware.RoleFinance,
		middleware.RoleDireksi,
	))
	karyawan.Get("/meta/departemen", karyawanCtrl.GetDepartemens)
	karyawan.Get("/meta/stats", karyawanCtrl.GetStats)
	karyawan.Get("/", karyawanCtrl.GetAll)
	karyawan.Get("/:id", karyawanCtrl.GetByID)
	karyawan.Post("/", middleware.RequireRole(middleware.RoleHRGA, middleware.RoleManager), karyawanCtrl.Create)
	karyawan.Put("/:id", middleware.RequireRole(middleware.RoleHRGA, middleware.RoleManager), karyawanCtrl.Update)
	karyawan.Delete("/:id", middleware.RequireRole(middleware.RoleHRGA), karyawanCtrl.Delete)
}
