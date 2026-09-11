package middleware

import (
	"sipadin/dto"

	"github.com/gofiber/fiber/v2"
)

// RequireRole memastikan user memiliki role yang diizinkan untuk mengakses route
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := GetClaims(ctx)

		for _, role := range allowedRoles {
			if claims.RoleSistem == role {
				return ctx.Next()
			}
		}

		return ctx.Status(fiber.StatusForbidden).JSON(dto.Fail("Akses ditolak: anda tidak memiliki izin untuk aksi ini"))
	}
}

// Konstanta role sesuai PRD v2.0
const (
	RoleKaryawan = "Karyawan"
	RoleManager  = "Manager"
	RoleHRGA     = "HRGA"
	RoleFinance  = "Finance"
	RoleDireksi  = "Direksi"
)
