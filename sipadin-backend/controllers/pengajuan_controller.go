package controllers

import (
	"fmt"
	"sipadin-backend/config"
	"sipadin-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GET: Ambil Semua Pengajuan atau filter per user_id
func GetPengajuan(c *fiber.Ctx) error {
	var listPengajuan []models.PengajuanDinas
	userID := c.Query("user_id")

	db := config.DB.Preload("User")
	if userID != "" && userID != "0" {
		db = db.Where("user_id = ?", userID)
	}

	db.Order("id desc").Find(&listPengajuan)
	return c.JSON(listPengajuan)
}

// GET: Ambil Detail 1 Pengajuan berdasarkan ID
func GetPengajuanByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var pengajuan models.PengajuanDinas

	if err := config.DB.Preload("User").First(&pengajuan, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Data pengajuan tidak ditemukan"})
	}

	return c.JSON(pengajuan)
}

// POST: Buat Pengajuan Baru (Karyawan)
func CreatePengajuan(c *fiber.Ctx) error {
	var pengajuan models.PengajuanDinas
	if err := c.BodyParser(&pengajuan); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Gagal membaca data input: " + err.Error()})
	}

	now := time.Now()
	if pengajuan.TanggalPengajuan == "" {
		pengajuan.TanggalPengajuan = now.Format("2006-01-02")
	}

	if pengajuan.NomorSurat == "" {
		var count int64
		config.DB.Model(&models.PengajuanDinas{}).Count(&count)
		deptCode := "IT"
		if pengajuan.Dept != "" {
			if len(pengajuan.Dept) >= 3 {
				deptCode = pengajuan.Dept[:3]
			}
		}
		pengajuan.NomorSurat = fmt.Sprintf("SPPD/%s/%d/%02d/%03d", deptCode, now.Year(), now.Month(), count+1)
	}

	if pengajuan.Status == "" {
		pengajuan.Status = "pending"
	}

	if err := config.DB.Create(&pengajuan).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan pengajuan: " + err.Error()})
	}

	// Preload User untuk respon lengkap
	config.DB.Preload("User").First(&pengajuan, pengajuan.ID)

	return c.Status(201).JSON(fiber.Map{
		"message": "Pengajuan berhasil dibuat",
		"data":    pengajuan,
	})
}

// PUT: Update Status Workflow (Atasan / HRGA / Notifikasi Direksi)
func UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	type StatusUpdateInput struct {
		Status                  string `json:"status"`
		Catatan                 string `json:"catatan"`
		CatatanAtasan           string `json:"catatan_atasan"`
		DisetujuiOleh           string `json:"disetujui_oleh"`
		TanggalDisetujui        string `json:"tanggal_disetujui"`
		CatatanHRGA             string `json:"catatan_hrga"`
		DiverifikasiOleh        string `json:"diverifikasi_oleh"`
		TanggalVerifikasiHRGA   string `json:"tanggal_verifikasi_hrga"`
		TanggalNotifikasiDireksi string `json:"tanggal_notifikasi_direksi"`
		PesanWA                 string `json:"pesan_wa"`
	}

	var input StatusUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input status tidak valid"})
	}

	var pengajuan models.PengajuanDinas
	if err := config.DB.First(&pengajuan, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Data pengajuan tidak ditemukan"})
	}

	todayStr := time.Now().Format("2006-01-02")

	if input.Status != "" {
		pengajuan.Status = input.Status
	}

	// Update field atasan jika status disetujui / ditolak atasan
	if input.Status == "approved" || input.Status == "rejected" {
		if input.CatatanAtasan != "" {
			pengajuan.CatatanAtasan = input.CatatanAtasan
		} else if input.Catatan != "" {
			pengajuan.CatatanAtasan = input.Catatan
		}
		if input.DisetujuiOleh != "" {
			pengajuan.DisetujuiOleh = input.DisetujuiOleh
		}
		if input.TanggalDisetujui != "" {
			pengajuan.TanggalDisetujui = input.TanggalDisetujui
		} else {
			pengajuan.TanggalDisetujui = todayStr
		}
	}

	// Update field HRGA jika verifikasi / tolak HRGA
	if input.Status == "verified_hrga" || input.Status == "rejected_hrga" {
		if input.CatatanHRGA != "" {
			pengajuan.CatatanHRGA = input.CatatanHRGA
		} else if input.Catatan != "" {
			pengajuan.CatatanHRGA = input.Catatan
		}
		if input.DiverifikasiOleh != "" {
			pengajuan.DiverifikasiOleh = input.DiverifikasiOleh
		}
		if input.TanggalVerifikasiHRGA != "" {
			pengajuan.TanggalVerifikasiHRGA = input.TanggalVerifikasiHRGA
		} else {
			pengajuan.TanggalVerifikasiHRGA = todayStr
		}
	}

	// Update notifikasi direksi
	if input.Status == "notified_direksi" {
		if input.TanggalNotifikasiDireksi != "" {
			pengajuan.TanggalNotifikasiDireksi = input.TanggalNotifikasiDireksi
		} else {
			pengajuan.TanggalNotifikasiDireksi = todayStr
		}
		if input.PesanWA != "" {
			pengajuan.PesanWA = input.PesanWA
		}

		notif := models.NotifikasiDireksi{
			PengajuanID:      pengajuan.ID,
			IsiPesan:         input.PesanWA,
			StatusPengiriman: "terkirim",
			SentAt:           time.Now(),
		}
		config.DB.Create(&notif)
	}

	if err := config.DB.Save(&pengajuan).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal memperbarui status pengajuan"})
	}

	config.DB.Preload("User").First(&pengajuan, pengajuan.ID)

	return c.JSON(fiber.Map{
		"message": "Status pengajuan berhasil diperbarui",
		"data":    pengajuan,
	})
}