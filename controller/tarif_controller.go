package controller

import (
	"sipadin/dto"
	"sipadin/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// TarifController menangani HTTP request untuk master tarif fasilitas
type TarifController struct {
	tarifSvc service.TarifService
}

func NewTarifController(tarifSvc service.TarifService) *TarifController {
	return &TarifController{tarifSvc: tarifSvc}
}

// GetAll - GET /api/tarif
func (c *TarifController) GetAll(ctx *fiber.Ctx) error {
	data, err := c.tarifSvc.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal mengambil data tarif"))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Daftar tarif fasilitas", data))
}

// GetByID - GET /api/tarif/:id
func (c *TarifController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}
	data, err := c.tarifSvc.GetByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data tarif ditemukan", data))
}

// Create - POST /api/tarif
func (c *TarifController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateTarifRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}
	data, err := c.tarifSvc.Create(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(dto.OK("Data tarif berhasil ditambahkan", data))
}

// Update - PUT /api/tarif/:id
func (c *TarifController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}
	var req dto.UpdateTarifRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}
	data, err := c.tarifSvc.Update(id, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data tarif berhasil diperbarui", data))
}

// Delete - DELETE /api/tarif/:id
func (c *TarifController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}
	if err := c.tarifSvc.Delete(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data tarif berhasil dihapus", nil))
}
