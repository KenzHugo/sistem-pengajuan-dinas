package controller

import (
	"strings"

	"sipadin/dto"
	"sipadin/service"

	"github.com/gofiber/fiber/v2"
)

// AuthController menangani HTTP request untuk autentikasi
type AuthController struct {
	authSvc service.AuthService
}

func NewAuthController(authSvc service.AuthService) *AuthController {
	return &AuthController{authSvc: authSvc}
}

// Login godoc
// POST /api/auth/login
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}

	req.Identifier = strings.TrimSpace(req.Identifier)
	req.Password = strings.TrimSpace(req.Password)

	if req.Identifier == "" || req.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Nomor ID dan password wajib diisi"))
	}

	result, err := c.authSvc.Login(req)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.Fail(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Login berhasil", result))
}

// GetDBUsers godoc
// GET /api/auth/users
func (c *AuthController) GetDBUsers(ctx *fiber.Ctx) error {
	users, err := c.authSvc.GetActiveUsers()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal mengambil data user dari database"))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data user berhasil diambil langsung dari database", users))
}
