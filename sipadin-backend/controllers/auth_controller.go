package controllers

import (
	"strings"

	"sipadin-backend/config"
	"sipadin-backend/models"

	"github.com/gofiber/fiber/v2"
)

// ============================================================
// POST /api/auth/login
// Body: { "nama": "...", "password": "..." }
//
// Lookup Karyawan_copy1 by Nama + Password + Aktif=1.
// Infers role dynamically:
//   - KodeDepartemen == "HRD" or "HRGA"   → "HRGA"
//   - KodeJabatan in atasan list           → "ATASAN"
//   - otherwise                            → "KARYAWAN"
// ============================================================

var atasanJabatanCodes = map[string]bool{
	"27": true, "28": true, "29": true, "30": true,
	"01": true, "02M": true, "04": true, "06": true, "08MS": true,
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Nama     string `json:"nama"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format input tidak valid",
		})
	}

	input.Nama     = strings.TrimSpace(input.Nama)
	input.Password = strings.TrimSpace(input.Password)

	if input.Nama == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nama dan kata sandi wajib diisi",
		})
	}

	// Lookup karyawan: Nama + Password + Aktif=1
	var karyawan models.KaryawanCopy1
	result := config.DB.
		Where("Nama = ? AND Password = ? AND Aktif = 1", input.Nama, input.Password).
		First(&karyawan)

	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Nama atau kata sandi salah, atau akun tidak aktif",
		})
	}

	role := inferRole(karyawan)

	userProfile := fiber.Map{
		"nomor_id":        karyawan.NomorID,
		"nama":            karyawan.Nama,
		"email":           karyawan.Email,
		"kode_departemen": karyawan.KodeDepartemen,
		"kode_jabatan":    karyawan.KodeJabatan,
		"kode_golongan":   karyawan.KodeGolongan,
		"role":            role,
		// Alias untuk kompatibilitas frontend
		"name":       karyawan.Nama,
		"nip":        karyawan.NomorID,
		"departemen": karyawan.KodeDepartemen,
		"jabatan":    karyawan.KodeJabatan,
		"avatar":     buildAvatar(karyawan.Nama),
	}

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"user":    userProfile,
	})
}

func inferRole(k models.KaryawanCopy1) string {
	dept := strings.ToUpper(strings.TrimSpace(k.KodeDepartemen))
	if dept == "HRD" || dept == "HRGA" {
		return "HRGA"
	}
	if atasanJabatanCodes[strings.TrimSpace(k.KodeJabatan)] {
		return "ATASAN"
	}
	return "KARYAWAN"
}

func buildAvatar(nama string) string {
	parts    := strings.Fields(nama)
	initials := ""
	for _, p := range parts {
		if len(p) > 0 {
			initials += string([]rune(p)[0])
		}
		if len(initials) >= 2 {
			break
		}
	}
	return strings.ToUpper(initials)
}