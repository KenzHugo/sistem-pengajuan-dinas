package controller

import (
	"sipadin/dto"
	"sipadin/middleware"
	"sipadin/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// SppdController menangani HTTP request untuk pengajuan SPPD
type SppdController struct {
	sppdSvc service.SppdService
}

func NewSppdController(sppdSvc service.SppdService) *SppdController {
	return &SppdController{sppdSvc: sppdSvc}
}

// Create - POST /api/sppd
func (c *SppdController) Create(ctx *fiber.Ctx) error {
	claims := middleware.GetClaims(ctx)
	var req dto.CreateSppdRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}

	result, err := c.sppdSvc.Create(claims.IDUser, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(dto.OK("Pengajuan berhasil disimpan sebagai draft", result))
}

// Submit - POST /api/sppd/:id/submit
func (c *SppdController) Submit(ctx *fiber.Ctx) error {
	claims := middleware.GetClaims(ctx)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}

	result, err := c.sppdSvc.Submit(id, claims.IDUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("SPPD berhasil diajukan ke Manager", result))
}

// GetByID - GET /api/sppd/:id
func (c *SppdController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}

	result, err := c.sppdSvc.GetByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data pengajuan ditemukan", result))
}

// GetMyList - GET /api/sppd/mine (Karyawan: lihat pengajuan sendiri)
func (c *SppdController) GetMyList(ctx *fiber.Ctx) error {
	claims := middleware.GetClaims(ctx)
	filter := dto.SppdListFilter{
		Status: ctx.Query("status"),
		Page:   parseInt(ctx.Query("page"), 1),
		Limit:  parseInt(ctx.Query("limit"), 20),
	}

	data, total, err := c.sppdSvc.GetMyList(claims.IDUser, filter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal mengambil data"))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.PaginatedResponse{
		Success: true, Message: "Data pengajuan saya",
		Data: data, Total: total,
		Page: filter.Page, Limit: filter.Limit,
	})
}

// GetAll - GET /api/sppd (semua role yang punya akses view)
func (c *SppdController) GetAll(ctx *fiber.Ctx) error {
	claims := middleware.GetClaims(ctx)
	filter := dto.SppdListFilter{
		Status: ctx.Query("status"),
		Page:   parseInt(ctx.Query("page"), 1),
		Limit:  parseInt(ctx.Query("limit"), 20),
	}

	var data []interface{}
	var total int64

	// Karyawan hanya lihat milik sendiri
	if claims.RoleSistem == "Karyawan" {
		sppds, t, err := c.sppdSvc.GetMyList(claims.IDUser, filter)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal mengambil data"))
		}
		for _, s := range sppds {
			data = append(data, s)
		}
		total = t
	} else {
		sppds, t, err := c.sppdSvc.GetAllList(filter)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal mengambil data"))
		}
		for _, s := range sppds {
			data = append(data, s)
		}
		total = t
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.PaginatedResponse{
		Success: true, Message: "Daftar pengajuan SPPD",
		Data: data, Total: total,
		Page: filter.Page, Limit: filter.Limit,
	})
}

// GetQueue - GET /api/sppd/queue?stage=manager|hrga|finance (untuk halaman approval)
func (c *SppdController) GetQueue(ctx *fiber.Ctx) error {
	stage := ctx.Query("stage")
	filter := dto.SppdListFilter{
		Page:  parseInt(ctx.Query("page"), 1),
		Limit: parseInt(ctx.Query("limit"), 20),
	}

	stageStatusMap := map[string][]string{
		"manager": {"pending_manager"},
		"hrga":    {"pending_hrga", "ready_to_send"},
		"finance": {"pending_finance"},
		"laporan": {"approved", "rejected"},
	}

	statuses, ok := stageStatusMap[stage]
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Stage tidak valid"))
	}

	data, total, err := c.sppdSvc.GetByStatus(statuses, filter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal mengambil data antrian"))
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.PaginatedResponse{
		Success: true, Message: "Antrian " + stage,
		Data: data, Total: total,
		Page: filter.Page, Limit: filter.Limit,
	})
}

// Update - PUT /api/sppd/:id
func (c *SppdController) Update(ctx *fiber.Ctx) error {
	claims := middleware.GetClaims(ctx)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}

	var req dto.UpdateSppdRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}

	result, err := c.sppdSvc.Update(id, claims.IDUser, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Pengajuan berhasil diperbarui", result))
}

// Delete - DELETE /api/sppd/:id
func (c *SppdController) Delete(ctx *fiber.Ctx) error {
	claims := middleware.GetClaims(ctx)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}

	if err := c.sppdSvc.Delete(id, claims.IDUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Pengajuan berhasil dihapus", nil))
}

func parseInt(s string, defaultVal int) int {
	v, err := strconv.Atoi(s)
	if err != nil || v <= 0 {
		return defaultVal
	}
	return v
}
