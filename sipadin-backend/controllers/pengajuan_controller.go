package controllers

import (
	"sipadin-backend/config"
	"sipadin-backend/models"

	"github.com/gofiber/fiber/v2"
)

// GET: Ambil Semua Pengajuan (Atasan/HRGA) atau per User (Karyawan)
func GetPengajuan(c *fiber.Ctx) error {
	var listPengajuan []models.PengajuanDinas
	userID := c.Query("user_id")

	db := config.DB.Preload("User")
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}

	db.Find(&listPengajuan)
	return c.JSON(listPengajuan)
}

// POST: Buat Pengajuan Baru (Karyawan)
func CreatePengajuan(c *fiber.Ctx) error {
	var pengajuan models.PengajuanDinas
	if err := c.BodyParser(&pengajuan); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Gagal membaca data input"})
	}

	pengajuan.Status = "pending" // Status awal selalu pending
	if err := config.DB.Create(&pengajuan).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan pengajuan"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Pengajuan berhasil dibuat",
		"data":    pengajuan,
	})
}

// PUT: Update Status Workflow (Atasan / HRGA)
func UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	type StatusUpdate struct {
		Status string `json:"status"`
	}

	var input StatusUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input status tidak valid"})
	}

	var pengajuan models.PengajuanDinas
	if err := config.DB.First(&pengajuan, id).Error; err != nil {
		return c.Status(444).JSON(fiber.Map{"message": "Data pengajuan tidak ditemukan"})
	}

	pengajuan.Status = input.Status
	config.DB.Save(&pengajuan)

	// Jika status berubah ke notified_direksi, buat log otomatis
	if input.Status == "notified_direksi" {
		notif := models.NotifikasiDireksi{
			PengajuanID: pengajuan.ID,
			IsiPesan:    "Notifikasi pengajuan dinas telah diteruskan ke direksi.",
		}
		config.DB.Create(&notif)
	}

	return c.JSON(fiber.Map{
		"message": "Status pengajuan berhasil diperbarui",
		"data":    pengajuan,
	})
}