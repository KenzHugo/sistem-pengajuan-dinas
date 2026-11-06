package middleware

import (
	"sipadin/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// RequireOperationalHours memastikan endpoint hanya dapat diakses pada jam kerja aktif (cross-check server & device)
func RequireOperationalHours(operasionalSvc service.OperasionalService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Jika user memiliki klaim spesial akses, lewatkan pembatasan
		if userToken, ok := ctx.Locals("user").(*jwt.Token); ok && userToken != nil {
			if claims, ok := userToken.Claims.(jwt.MapClaims); ok {
				if isSpecial, ok := claims["is_special_access"].(bool); ok && isSpecial {
					return ctx.Next()
				}
			}
		}

		clientTime := ctx.Get("X-Client-Time")
		if clientTime == "" {
			clientTime = ctx.Query("client_time")
		}

		allowed, msg, _ := operasionalSvc.ValidateAccess("", clientTime)
		if !allowed {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success":            false,
				"message":            msg,
				"operational_closed": true,
			})
		}

		return ctx.Next()
	}
}
