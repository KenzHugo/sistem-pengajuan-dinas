package controller

import (
	"encoding/csv"
	"fmt"
	"sipadin/config"
	"sipadin/dto"
	"sipadin/middleware"
	"sipadin/models"
	"sipadin/service"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

// GetSppdLedger - GET /api/finance/sppd-ledger/:id
func (c *ApprovalController) GetSppdLedger(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail("ID SPPD tidak valid"))
	}

	result, err := c.approvalSvc.GetSppdLedger(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.Fail(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(dto.OK("Data ledger keuangan berhasil dimuat", result))
}

// ExportFinancialCSV - GET /api/finance/export-csv
// Mengekspor rincian keuangan SPPD per personil (base query: sppd_detail_biaya) menggunakan encoding/csv
func (c *ApprovalController) ExportFinancialCSV(ctx *fiber.Ctx) error {
	status := ctx.Query("status", "")

	query := config.DB.Model(&models.SppdPengajuan{}).
		Preload("Pemohon").
		Preload("Anggota").
		Preload("Anggota.User").
		Preload("DetailBiaya", func(db *gorm.DB) *gorm.DB {
			return db.Order("id_detail ASC")
		}).
		Order("id_sppd DESC")

	if status != "" {
		query = query.Where("status_dokumen = ?", status)
	}

	var sppdList []models.SppdPengajuan
	if err := query.Find(&sppdList).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.Fail("Gagal memuat data laporan"))
	}

	ctx.Set("Content-Type", "text/csv; charset=utf-8")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=Laporan_Finance_SPPD_Per_Personil_%s.csv", time.Now().Format("20060102_150405")))

	// Tulis UTF-8 BOM agar Excel menampilkan aksen dan teks dengan rapi
	ctx.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(ctx.Response().BodyWriter())
	defer writer.Flush()

	// 1. Header CSV Per Kepala
	header := []string{
		"No SPPD",
		"Tanggal Berangkat",
		"Tanggal Kembali",
		"Tujuan Dinas",
		"Deskripsi Kegiatan",
		"Nama Personil",
		"Role",
		"Jabatan / Kategori",
		"Uang Makan (Rp)",
		"Uang Saku (Rp)",
		"Uang Hotel (Rp)",
		"Subtotal Hak Personil (Rp)",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// 2. Loop SPPD dan rincikan per kepala
	for _, sppd := range sppdList {
		var details []models.SppdDetailBiaya
		// Ambil ledger detail terverifikasi dan tervalidasi
		ledgerRes, _ := c.approvalSvc.GetSppdLedger(sppd.IDSppd)
		if ledgerRes != nil && len(ledgerRes.DetailLedger) > 0 {
			details = ledgerRes.DetailLedger
		}

		tujuan := sppd.KotaTujuan
		if sppd.InstansiTujuan != "" {
			tujuan = fmt.Sprintf("%s (%s)", sppd.KotaTujuan, sppd.InstansiTujuan)
		}
		tglB := sppd.TglBerangkat.Format("2006-01-02")
		tglK := sppd.TglKembali.Format("2006-01-02")

		// Rincian masing-masing personil dalam SPPD
		for _, d := range details {
			row := []string{
				sppd.NoSppd,
				tglB,
				tglK,
				tujuan,
				d.SnapshotDeskripsi,
				d.SnapshotNama,
				d.RolePersonil,
				d.SnapshotJabatan,
				fmt.Sprintf("%.2f", d.UangMakan),
				fmt.Sprintf("%.2f", d.UangSaku),
				fmt.Sprintf("%.2f", d.UangHotel),
				fmt.Sprintf("%.2f", d.SubtotalHakIndividu),
			}
			_ = writer.Write(row)
		}

		// 3. Tambahkan 1 baris ekstra untuk Biaya Transportasi Bersama / Tiket (jika ada)
		biayaBersama := sppd.BiayaTransport + sppd.BiayaTiket
		if biayaBersama > 0 {
			ketBersama := "Transportasi Bersama"
			if sppd.BiayaTransport > 0 && sppd.BiayaTiket > 0 {
				ketBersama = fmt.Sprintf("Transport Bersama (Rp%.0f) + Tiket (Rp%.0f)", sppd.BiayaTransport, sppd.BiayaTiket)
			} else if sppd.BiayaTiket > 0 {
				ketBersama = "Tiket Bersama"
			}

			rowBersama := []string{
				sppd.NoSppd,
				tglB,
				tglK,
				tujuan,
				sppd.MaksudTujuan,
				fmt.Sprintf("[BIAYA BERSAMA] %s", ketBersama),
				"Biaya Bersama",
				sppd.Transportasi,
				"0.00",
				"0.00",
				"0.00",
				fmt.Sprintf("%.2f", biayaBersama),
			}
			_ = writer.Write(rowBersama)
		}
	}

	return nil
}


