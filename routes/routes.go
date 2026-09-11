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

	// ─── Tarif Routes (HRGA Full Access, semua lainnya GET only) ──────────────
	tarif := protected.Group("/tarif")
	tarif.Get("/", tarifCtrl.GetAll)
	tarif.Get("/:id", tarifCtrl.GetByID)
	tarif.Post("/", middleware.RequireRole(middleware.RoleHRGA), tarifCtrl.Create)
	tarif.Put("/:id", middleware.RequireRole(middleware.RoleHRGA), tarifCtrl.Update)
	tarif.Delete("/:id", middleware.RequireRole(middleware.RoleHRGA), tarifCtrl.Delete)

	// ─── User Management Routes (HRGA Full Access) ────────────────────────────
	users := protected.Group("/users")
	users.Get("/", userCtrl.GetAll) // semua bisa lihat (untuk select anggota)
	users.Get("/list", userCtrl.GetAll)
	users.Get("/:id", userCtrl.GetByID)
	users.Post("/", middleware.RequireRole(middleware.RoleHRGA), userCtrl.Create)
	users.Put("/:id", middleware.RequireRole(middleware.RoleHRGA), userCtrl.Update)
	users.Delete("/:id", middleware.RequireRole(middleware.RoleHRGA), userCtrl.Delete)
}
