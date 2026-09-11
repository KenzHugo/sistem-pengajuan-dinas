package controller

import (
	"sipadin/dto"
	"sipadin/middleware"
	"sipadin/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ApprovalController menangani HTTP request untuk workflow approval
type ApprovalController struct {
	approvalSvc service.ApprovalService
}

func NewApprovalController(approvalSvc service.ApprovalService) *ApprovalController {
	return &ApprovalController{approvalSvc: approvalSvc}
}

// ManagerAction - POST /api/approval/:id/manager
func (c *ApprovalController) ManagerAction(ctx *fiber.Ctx) error {
	claims := middleware.GetClaims(ctx)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}

	var req dto.ApprovalRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}

	result, err := c.approvalSvc.ManagerAction(id, claims.IDUser, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Tindakan Manager berhasil", result))
}

// HRGAAction - POST /api/approval/:id/hrga
func (c *ApprovalController) HRGAAction(ctx *fiber.Ctx) error {
	claims := middleware.GetClaims(ctx)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}

	var req dto.ApprovalRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}

	result, err := c.approvalSvc.HRGAAction(id, claims.IDUser, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Tindakan HRGA berhasil", result))
}

// FinanceAction - POST /api/approval/:id/finance
func (c *ApprovalController) FinanceAction(ctx *fiber.Ctx) error {
	claims := middleware.GetClaims(ctx)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}

	var req dto.ApprovalRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("Format request tidak valid"))
	}

	result, err := c.approvalSvc.FinanceAction(id, claims.IDUser, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Tindakan Finance berhasil", result))
}

// SendToDirector - POST /api/approval/:id/send
func (c *ApprovalController) SendToDirector(ctx *fiber.Ctx) error {
	claims := middleware.GetClaims(ctx)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID tidak valid"))
	}

	result, err := c.approvalSvc.SendToDirector(id, claims.IDUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Dokumen berhasil dikirim ke Direksi via WhatsApp", result))
}
