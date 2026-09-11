package middleware

import (
	"sipadin/config"
	"sipadin/dto"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected adalah middleware yang memvalidasi JWT token
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.AppConfig.JWTSecret)},
		ErrorHandler: jwtErrorHandler,
	})
}

func jwtErrorHandler(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(dto.Fail("Token tidak valid atau sudah kadaluarsa"))
}

// UserClaims adalah struktur data yang diekstrak dari JWT
type UserClaims struct {
	IDUser      int
	NamaLengkap string
	RoleSistem  string
}

// GetClaims mengekstrak data user dari JWT token yang sudah divalidasi
func GetClaims(ctx *fiber.Ctx) *UserClaims {
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	idUser := 0
	if v, ok := claims["id_user"].(float64); ok {
		idUser = int(v)
	}

	return &UserClaims{
		IDUser:      idUser,
		NamaLengkap: claims["nama_lengkap"].(string),
		RoleSistem:  claims["role_sistem"].(string),
	}
}
