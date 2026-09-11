package controller

import (
	"sipadin/dto"
	"sipadin/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// UserController menangani HTTP request untuk manajemen user
type UserController struct {
	userSvc service.UserService
}

func NewUserController(userSvc service.UserService) *UserController {
	return &UserController{userSvc: userSvc}
}

// GetAll - GET /api/users
func (c *UserController) GetAll(ctx *fiber.Ctx) error {
	data, err := c.userSvc.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal mengambil data user"))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Daftar user", data))
}

// GetByID - GET /api/users/:id
func (c *UserController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}
	data, err := c.userSvc.GetByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data user ditemukan", data))
}

// Create - POST /api/users
func (c *UserController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}
	data, err := c.userSvc.Create(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(dto.OK("User berhasil dibuat", data))
}

// Update - PUT /api/users/:id
func (c *UserController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}
	var req dto.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}
	data, err := c.userSvc.Update(id, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("User berhasil diperbarui", data))
}

// Delete - DELETE /api/users/:id
func (c *UserController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}
	if err := c.userSvc.Delete(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("User berhasil dihapus", nil))
}
