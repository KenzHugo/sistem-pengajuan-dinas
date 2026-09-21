package controller

import (
	"sipadin/dto"
	"sipadin/service"

	"github.com/gofiber/fiber/v2"
)

// TarifController menangani HTTP request untuk master tarif perjalanan dinas
type TarifController struct {
	tarifSvc service.TarifService
}

func NewTarifController(tarifSvc service.TarifService) *TarifController {
	return &TarifController{tarifSvc: tarifSvc}
}

// GetAll - GET /api/tarif
func (c *TarifController) GetAll(ctx *fiber.Ctx) error {
	area := ctx.Query("area", "")
	tarifs, err := c.tarifSvc.GetAll(area)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal mengambil data tarif: " + err.Error()))
	}

	fasilitas, _ := c.tarifSvc.GetFasilitas()
	areas, _ := c.tarifSvc.GetAreas()

	responsePayload := fiber.Map{
		"tarifs":    tarifs,
		"fasilitas": fasilitas,
		"areas":     areas,
		"total":     len(tarifs),
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Daftar master tarif perjalanandinas_tarif", responsePayload))
}

// GetFasilitas - GET /api/tarif/fasilitas
func (c *TarifController) GetFasilitas(ctx *fiber.Ctx) error {
	fasilitas, err := c.tarifSvc.GetFasilitas()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal mengambil data fasilitas: " + err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Daftar fasilitas perjalanan dinas", fasilitas))
}

// GetAreas - GET /api/tarif/areas
func (c *TarifController) GetAreas(ctx *fiber.Ctx) error {
	areas, err := c.tarifSvc.GetAreas()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal mengambil daftar area: " + err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Daftar area tarif", areas))
}

// CreateOrUpdate - POST /api/tarif
func (c *TarifController) Create(ctx *fiber.Ctx) error {
	var req dto.TarifRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}
	data, err := c.tarifSvc.CreateOrUpdate(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data tarif berhasil disimpan", data))
}

// Update - PUT /api/tarif
func (c *TarifController) Update(ctx *fiber.Ctx) error {
	var req dto.TarifRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}
	data, err := c.tarifSvc.CreateOrUpdate(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data tarif berhasil diperbarui", data))
}

// Delete - DELETE /api/tarif
func (c *TarifController) Delete(ctx *fiber.Ctx) error {
	area := ctx.Query("area")
	jabatan := ctx.Query("jabatan")
	if area == "" || jabatan == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Parameter area dan jabatan wajib diisi"))
	}
	if err := c.tarifSvc.Delete(area, jabatan); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data tarif berhasil dihapus", nil))
}

// GetByID untuk backward compatibility
func (c *TarifController) GetByID(ctx *fiber.Ctx) error {
	return c.GetAll(ctx)
}
