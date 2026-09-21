package controller

import (
	"sipadin/dto"
	"sipadin/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// KaryawanController menangani HTTP request untuk data karyawan
type KaryawanController struct {
	karyawanSvc service.KaryawanService
}

// NewKaryawanController membuat instance baru KaryawanController
func NewKaryawanController(karyawanSvc service.KaryawanService) *KaryawanController {
	return &KaryawanController{karyawanSvc: karyawanSvc}
}

// GetAll - GET /api/karyawan
func (c *KaryawanController) GetAll(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "15"))

	query := dto.KaryawanQuery{
		Search:         ctx.Query("search"),
		KodeDepartemen: ctx.Query("kode_departemen"),
		Aktif:          ctx.Query("aktif"),
		Page:           page,
		Limit:          limit,
		SortBy:         ctx.Query("sort_by"),
		Order:          ctx.Query("order"),
	}

	data, err := c.karyawanSvc.GetAll(query)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal mengambil data karyawan: " + err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Daftar karyawan", data))
}

// GetByID - GET /api/karyawan/:id
func (c *KaryawanController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Nomor ID tidak valid"))
	}

	data, err := c.karyawanSvc.GetByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(dto.Fail(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data karyawan ditemukan", data))
}

// Create - POST /api/karyawan
func (c *KaryawanController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateKaryawanRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}

	data, err := c.karyawanSvc.Create(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.OK("Data karyawan berhasil ditambahkan", data))
}

// Update - PUT /api/karyawan/:id
func (c *KaryawanController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Nomor ID tidak valid"))
	}

	var req dto.UpdateKaryawanRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}

	data, err := c.karyawanSvc.Update(id, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data karyawan berhasil diperbarui", data))
}

// Delete - DELETE /api/karyawan/:id
func (c *KaryawanController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Nomor ID tidak valid"))
	}

	if err := c.karyawanSvc.Delete(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data karyawan berhasil dihapus", nil))
}

// GetDepartemens - GET /api/karyawan/meta/departemen
func (c *KaryawanController) GetDepartemens(ctx *fiber.Ctx) error {
	depts, err := c.karyawanSvc.GetDepartemenList()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal memuat daftar departemen"))
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Daftar departemen", depts))
}

// GetStats - GET /api/karyawan/meta/stats
func (c *KaryawanController) GetStats(ctx *fiber.Ctx) error {
	stats, err := c.karyawanSvc.GetStats()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal memuat statistik karyawan"))
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Statistik karyawan", stats))
}
