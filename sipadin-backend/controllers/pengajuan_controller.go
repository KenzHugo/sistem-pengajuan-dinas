package controllers

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"sipadin-backend/config"
	"sipadin-backend/models"

	"github.com/gofiber/fiber/v2"
)

// ============================================================
// STATUS FLOW CONSTANTS
// ============================================================
const (
	StatusSubmitted        = "SUBMITTED"
	StatusApprovedAtasan   = "APPROVED_ATASAN"
	StatusControlledHRD    = "CONTROLLED_HRD"
	StatusWASentDireksi    = "WA_SENT_DIREKSI"
	StatusConfirmedDireksi = "CONFIRMED_DIREKSI"
	StatusSuratTugasIssued = "SURAT_TUGAS_ISSUED"
	StatusRejected         = "REJECTED"
)

// ============================================================
// HELPER: enrich pengajuan with Karyawan & Atasan details
// ============================================================
func enrichPengajuan(p *models.PerjalananDinasPengajuan) {
	if p.NomorID != "" {
		var k models.KaryawanCopy1
		if err := config.DB.Where("NomorID = ?", p.NomorID).First(&k).Error; err == nil {
			p.Karyawan = &k
			p.Nama = k.Nama
			p.KodeDepartemen = k.KodeDepartemen
			p.KodeJabatan = k.KodeJabatan
		}
	}
	if p.AtasanNomorID != "" {
		var a models.KaryawanCopy1
		if err := config.DB.Where("NomorID = ?", p.AtasanNomorID).First(&a).Error; err == nil {
			p.Atasan = &a
			p.AtasanNama = a.Nama
		}
	}
}

// ============================================================
// HELPER: extract Pengajuan ID / NoPengajuan from params or query
// ============================================================
func extractPengajuanID(c *fiber.Ctx) string {
	id := c.Query("id")
	if id == "" {
		id = c.Query("no")
	}
	if id == "" {
		id = c.Query("no_pengajuan")
	}
	if id == "" {
		id = c.Params("*")
	}
	if id == "" {
		id = c.Params("id")
	}
	if unescaped, err := url.QueryUnescape(id); err == nil && unescaped != "" {
		id = unescaped
	}
	return strings.TrimSpace(id)
}

// ============================================================
// HELPER: fetch pengajuan by NoPengajuan + preload Karyawan
// ============================================================
func getPengajuanWithKaryawan(noPengajuan string) (*models.PerjalananDinasPengajuan, error) {
	noPengajuan = strings.TrimSpace(noPengajuan)
	if unescaped, err := url.QueryUnescape(noPengajuan); err == nil && unescaped != "" {
		noPengajuan = unescaped
	}
	var p models.PerjalananDinasPengajuan
	err := config.DB.Where("NoPengajuan = ? OR NomorSuratTugas = ?", noPengajuan, noPengajuan).First(&p).Error
	if err != nil {
		err = config.DB.Where("NoPengajuan LIKE ?", "%"+noPengajuan).First(&p).Error
		if err != nil {
			return nil, err
		}
	}
	enrichPengajuan(&p)
	return &p, nil
}

// ============================================================
// HELPER: generate NoPengajuan — format SPD/YYYYMM/XXXX
// ============================================================
func generateNoPengajuan() (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("SPD/%04d%02d/", now.Year(), int(now.Month()))

	var count int64
	config.DB.Model(&models.PerjalananDinasPengajuan{}).
		Where("NoPengajuan LIKE ?", prefix+"%").
		Count(&count)

	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}

// ============================================================
// HELPER: generate NomorSuratTugas — format ST/YYYY/MM/XXXX
// ============================================================
func generateNomorSuratTugas() (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("ST/%04d/%02d/", now.Year(), int(now.Month()))

	var count int64
	config.DB.Model(&models.PerjalananDinasPengajuan{}).
		Where("NomorSuratTugas LIKE ?", prefix+"%").
		Count(&count)

	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}

// ============================================================
// GET /api/master/tarif
// Returns all PerjalananDinas_Tarif records
// ============================================================
func GetMasterTarif(c *fiber.Ctx) error {
	var tarif []models.PerjalananDinasTarif
	if err := config.DB.Order("Area, NoUrut").Find(&tarif).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data tarif: " + err.Error(),
		})
	}
	return c.JSON(tarif)
}

// ============================================================
// GET /api/master/fasilitas
// Returns all PerjalananDinas_Fasilitas records
// ============================================================
func GetMasterFasilitas(c *fiber.Ctx) error {
	var fasilitas []models.PerjalananDinasFasilitas
	if err := config.DB.Find(&fasilitas).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data fasilitas: " + err.Error(),
		})
	}
	return c.JSON(fasilitas)
}

// ============================================================
// POST /api/pengajuan
// Karyawan mengajukan perjalanan dinas baru.
// Auto-generates NoPengajuan and sets StatusFlow = SUBMITTED.
// ============================================================
func CreatePengajuan(c *fiber.Ctx) error {
	type CreateInput struct {
		NomorID           string  `json:"nomor_id"`
		AtasanNomorID     string  `json:"atasan_nomor_id"`
		AreaTujuan        string  `json:"area_tujuan"`
		KotaTujuan        string  `json:"kota_tujuan"`
		MaksudTujuan      string  `json:"maksud_tujuan"`
		KriteriaFasilitas string  `json:"kriteria_fasilitas"`
		TanggalBerangkat  string  `json:"tanggal_berangkat"` // "YYYY-MM-DD"
		TanggalKembali    string  `json:"tanggal_kembali"`   // "YYYY-MM-DD"
		JamBerangkat      string  `json:"jam_berangkat"`
		JamKembali        string  `json:"jam_kembali"`
		EstimasiBiaya     float64 `json:"estimasi_biaya"`
	}

	var input CreateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format input tidak valid: " + err.Error(),
		})
	}

	if input.NomorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "NomorID wajib diisi",
		})
	}

	// Parse dates
	tglBerangkat, err := time.Parse("2006-01-02", input.TanggalBerangkat)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format tanggal_berangkat tidak valid (gunakan YYYY-MM-DD)",
		})
	}
	tglKembali, err := time.Parse("2006-01-02", input.TanggalKembali)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format tanggal_kembali tidak valid (gunakan YYYY-MM-DD)",
		})
	}

	noPengajuan, err := generateNoPengajuan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal generate nomor pengajuan",
		})
	}

	if input.JamBerangkat == "" {
		input.JamBerangkat = "08:00"
	}
	if input.JamKembali == "" {
		input.JamKembali = "17:00"
	}
	if input.KriteriaFasilitas == "" {
		input.KriteriaFasilitas = "1. < 200 KM & < 8 Jam"
	}

	pengajuan := models.PerjalananDinasPengajuan{
		NoPengajuan:       noPengajuan,
		NomorID:           input.NomorID,
		AtasanNomorID:     input.AtasanNomorID,
		AreaTujuan:        input.AreaTujuan,
		KotaTujuan:        input.KotaTujuan,
		MaksudTujuan:      input.MaksudTujuan,
		KriteriaFasilitas: input.KriteriaFasilitas,
		TanggalBerangkat:  tglBerangkat,
		TanggalKembali:    tglKembali,
		JamBerangkat:      input.JamBerangkat,
		JamKembali:        input.JamKembali,
		EstimasiBiaya:     input.EstimasiBiaya,
		StatusFlow:        StatusSubmitted,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := config.DB.Create(&pengajuan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan pengajuan: " + err.Error(),
		})
	}

	// Enrich with karyawan info for response
	result, _ := getPengajuanWithKaryawan(noPengajuan)
	if result == nil {
		result = &pengajuan
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pengajuan berhasil dibuat",
		"data":    result,
	})
}

// ============================================================
// GET /api/pengajuan
// Query params:
//   role     = KARYAWAN | ATASAN | HRGA
//   user_id  = NomorID of the requesting user
//   status   = (optional) filter by StatusFlow (or 'ALL'/'SEMUA' for all)
// ============================================================
func GetPengajuan(c *fiber.Ctx) error {
	role := strings.ToUpper(strings.TrimSpace(c.Query("role")))
	userID := strings.TrimSpace(c.Query("user_id"))
	statusFilter := strings.TrimSpace(c.Query("status"))
	if strings.ToUpper(statusFilter) == "ALL" || strings.ToUpper(statusFilter) == "SEMUA" {
		statusFilter = ""
	}

	db := config.DB.Model(&models.PerjalananDinasPengajuan{})

	switch role {
	case "KARYAWAN":
		// Karyawan sees their own submissions
		if userID != "" {
			db = db.Where("NomorID = ?", userID)
		}
		if statusFilter != "" {
			db = db.Where("StatusFlow = ?", statusFilter)
		}
	case "ATASAN":
		// Atasan: when reviewing submissions, show SUBMITTED records from subordinates
		// or records assigned to them. If statusFilter is provided, honor it.
		if statusFilter != "" {
			if userID != "" {
				db = db.Where("(AtasanNomorID = ? OR AtasanNomorID = '' OR AtasanNomorID IS NULL) AND StatusFlow = ?", userID, statusFilter)
			} else {
				db = db.Where("StatusFlow = ?", statusFilter)
			}
		} else {
			// Default overview for Atasan: show all SUBMITTED (needing decision) + anything assigned to them
			if userID != "" {
				db = db.Where("AtasanNomorID = ? OR AtasanNomorID = '' OR AtasanNomorID IS NULL OR StatusFlow = ?", userID, StatusSubmitted)
			}
		}
	case "HRGA":
		// HRGA sees all records (or filtered by status if provided)
		if statusFilter != "" {
			db = db.Where("StatusFlow = ?", statusFilter)
		}
	default:
		// General / reporting query
		if statusFilter != "" {
			db = db.Where("StatusFlow = ?", statusFilter)
		}
		if userID != "" {
			db = db.Where("NomorID = ?", userID)
		}
	}

	var list []models.PerjalananDinasPengajuan
	if err := db.Order("CreatedAt DESC").Find(&list).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pengajuan: " + err.Error(),
		})
	}

	// Enrich each record with full karyawan & atasan info
	for i := range list {
		enrichPengajuan(&list[i])
	}

	return c.JSON(list)
}

// ============================================================
// GET /api/pengajuan/:id
// Returns detail of a single pengajuan by NoPengajuan
// ============================================================
func GetPengajuanByID(c *fiber.Ctx) error {
	id := extractPengajuanID(c)
	p, err := getPengajuanWithKaryawan(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pengajuan tidak ditemukan",
		})
	}
	return c.JSON(p)
}

// ============================================================
// PUT /api/pengajuan/:id/atasan-approve
// Atasan menyetujui pengajuan (SUBMITTED → APPROVED_ATASAN)
// ============================================================
func AtasanApprove(c *fiber.Ctx) error {
	id := extractPengajuanID(c)
	p, err := getPengajuanWithKaryawan(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Pengajuan tidak ditemukan"})
	}
	if p.StatusFlow != StatusSubmitted {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": fmt.Sprintf("Tidak dapat menyetujui: status saat ini adalah '%s', bukan '%s'", p.StatusFlow, StatusSubmitted),
		})
	}

	type Input struct {
		AtasanNomorID string `json:"atasan_nomor_id"`
		Catatan       string `json:"catatan"`
	}
	var input Input
	_ = c.BodyParser(&input)
	if input.AtasanNomorID != "" {
		p.AtasanNomorID = input.AtasanNomorID
	}
	if input.Catatan != "" {
		p.CatatanHRD = input.Catatan
	}

	p.StatusFlow = StatusApprovedAtasan
	p.UpdatedAt = time.Now()

	if err := config.DB.Save(p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui status: " + err.Error()})
	}
	enrichPengajuan(p)
	return c.JSON(fiber.Map{"message": "Pengajuan disetujui oleh Atasan", "data": p})
}

// ============================================================
// PUT /api/pengajuan/:id/atasan-reject / PUT /api/pengajuan/:id/reject
// Atasan / HRGA / Sistem menolak pengajuan
// ============================================================
func AtasanReject(c *fiber.Ctx) error {
	id := extractPengajuanID(c)
	p, err := getPengajuanWithKaryawan(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Pengajuan tidak ditemukan"})
	}
	if p.StatusFlow == StatusSuratTugasIssued || p.StatusFlow == StatusRejected {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": fmt.Sprintf("Tidak dapat menolak: status saat ini sudah '%s'", p.StatusFlow),
		})
	}

	type Input struct {
		CatatanHRD string `json:"catatan_hrd"`
		Catatan    string `json:"catatan"`
		HRDNomorID string `json:"hrd_nomor_id"`
	}
	var input Input
	_ = c.BodyParser(&input)
	if input.CatatanHRD != "" {
		p.CatatanHRD = input.CatatanHRD
	} else if input.Catatan != "" {
		p.CatatanHRD = input.Catatan
	}
	if input.HRDNomorID != "" {
		p.HRDNomorID = input.HRDNomorID
	}

	p.StatusFlow = StatusRejected
	p.UpdatedAt = time.Now()

	if err := config.DB.Save(p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui status: " + err.Error()})
	}
	enrichPengajuan(p)
	return c.JSON(fiber.Map{"message": "Pengajuan ditolak", "data": p})
}

// ============================================================
// PUT /api/pengajuan/:id/hrd-reject / PUT /api/pengajuan/action/hrd-reject
// HRGA menolak pengajuan di tahap kontrol kebijakan (APPROVED_ATASAN / CONTROLLED_HRD → REJECTED)
// ============================================================
func HRDReject(c *fiber.Ctx) error {
	id := extractPengajuanID(c)
	p, err := getPengajuanWithKaryawan(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Pengajuan tidak ditemukan"})
	}
	if p.StatusFlow == StatusSuratTugasIssued || p.StatusFlow == StatusRejected {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": fmt.Sprintf("Tidak dapat menolak: status saat ini adalah '%s'", p.StatusFlow),
		})
	}

	type Input struct {
		CatatanHRD string `json:"catatan_hrd"`
		Catatan    string `json:"catatan"`
		HRDNomorID string `json:"hrd_nomor_id"`
	}
	var input Input
	_ = c.BodyParser(&input)

	note := input.CatatanHRD
	if note == "" {
		note = input.Catatan
	}
	if note == "" {
		note = "Ditolak oleh HRGA / Tidak memenuhi kebijakan anggaran"
	}

	p.CatatanHRD = note
	if input.HRDNomorID != "" {
		p.HRDNomorID = input.HRDNomorID
	}
	p.StatusFlow = StatusRejected
	p.UpdatedAt = time.Now()

	if err := config.DB.Save(p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui status penolakan HRGA: " + err.Error()})
	}
	enrichPengajuan(p)
	return c.JSON(fiber.Map{"message": "Pengajuan berhasil ditolak oleh HRGA", "data": p})
}

// ============================================================
// PUT /api/pengajuan/:id/direksi-reject / PUT /api/pengajuan/action/direksi-reject
// Direksi / HRGA menginput penolakan Direksi (WA_SENT_DIREKSI / CONTROLLED_HRD → REJECTED)
// ============================================================
func DireksiReject(c *fiber.Ctx) error {
	id := extractPengajuanID(c)
	p, err := getPengajuanWithKaryawan(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Pengajuan tidak ditemukan"})
	}
	if p.StatusFlow == StatusSuratTugasIssued || p.StatusFlow == StatusRejected {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": fmt.Sprintf("Tidak dapat menolak: status saat ini adalah '%s'", p.StatusFlow),
		})
	}

	type Input struct {
		KonfirmasiNote string `json:"konfirmasi_note"`
		Catatan        string `json:"catatan"`
		CatatanHRD     string `json:"catatan_hrd"`
	}
	var input Input
	_ = c.BodyParser(&input)

	note := input.KonfirmasiNote
	if note == "" {
		note = input.Catatan
	}
	if note == "" {
		note = input.CatatanHRD
	}
	if note == "" {
		note = "Ditolak oleh Direksi"
	}

	now := time.Now()
	p.KonfirmasiDireksiNote = note
	p.TanggalKonfirmasiDireksi = &now
	p.CatatanHRD = note
	p.StatusFlow = StatusRejected
	p.UpdatedAt = now

	if err := config.DB.Save(p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui status penolakan Direksi: " + err.Error()})
	}
	enrichPengajuan(p)
	return c.JSON(fiber.Map{"message": "Pengajuan berhasil ditolak oleh Direksi", "data": p})
}


// ============================================================
// PUT /api/pengajuan/:id/hrd-control
// HRGA memverifikasi & mengontrol (APPROVED_ATASAN → CONTROLLED_HRD)
// Body: { "catatan_hrd": "...", "estimasi_biaya": 0.0, "hrd_nomor_id": "..." }
// ============================================================
func HRDControl(c *fiber.Ctx) error {
	id := extractPengajuanID(c)
	p, err := getPengajuanWithKaryawan(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Pengajuan tidak ditemukan"})
	}
	if p.StatusFlow != StatusApprovedAtasan {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": fmt.Sprintf("Tidak dapat dikontrol HRD: status saat ini adalah '%s', bukan '%s'", p.StatusFlow, StatusApprovedAtasan),
		})
	}

	type Input struct {
		CatatanHRD    string  `json:"catatan_hrd"`
		EstimasiBiaya float64 `json:"estimasi_biaya"`
		HRDNomorID    string  `json:"hrd_nomor_id"`
	}
	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Format input tidak valid"})
	}

	if input.CatatanHRD != "" {
		p.CatatanHRD = input.CatatanHRD
	}
	if input.EstimasiBiaya > 0 {
		p.EstimasiBiaya = input.EstimasiBiaya
	}
	if input.HRDNomorID != "" {
		p.HRDNomorID = input.HRDNomorID
	}

	p.StatusFlow = StatusControlledHRD
	p.UpdatedAt = time.Now()

	if err := config.DB.Save(p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui status: " + err.Error()})
	}
	enrichPengajuan(p)
	return c.JSON(fiber.Map{"message": "Pengajuan berhasil dikontrol oleh HRD/HRGA", "data": p})
}

// ============================================================
// GET /api/pengajuan/:id/wa-text
// HRGA mempersiapkan pesan WhatsApp ke Direksi.
// Updates StatusFlow → WA_SENT_DIREKSI.
// Returns: { text: "...", whatsapp_url: "https://api.whatsapp.com/send?text=..." }
// ============================================================
func GetWAText(c *fiber.Ctx) error {
	id := extractPengajuanID(c)
	p, err := getPengajuanWithKaryawan(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Pengajuan tidak ditemukan"})
	}
	if p.StatusFlow != StatusControlledHRD && p.StatusFlow != StatusWASentDireksi {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": fmt.Sprintf("Status harus '%s' untuk mengirim WA Direksi, saat ini: '%s'", StatusControlledHRD, p.StatusFlow),
		})
	}

	// Build employee info
	nama := p.NomorID
	kodeDept := "-"
	if p.Karyawan != nil {
		nama = p.Karyawan.Nama
		kodeDept = p.Karyawan.KodeDepartemen
	}

	tglBerangkat := p.TanggalBerangkat.Format("02 Jan 2006")
	tglKembali := p.TanggalKembali.Format("02 Jan 2006")

	waText := fmt.Sprintf(
		"Bapak/Ibu Direksi, izin menginformasikan rencana perjalanan dinas atas nama %s (%s), "+
			"Departemen %s, tujuan %s (%s) pada tanggal %s s.d. %s untuk keperluan %s. "+
			"Estimasi biaya: Rp %.0f. "+
			"Pengajuan telah diverifikasi oleh atasan dan dikontrol oleh HRD/HRGA. "+
			"Mohon konfirmasi persetujuan: Balas 'OK / Disetujui'. Terima kasih.",
		nama, p.NomorID,
		kodeDept,
		p.KotaTujuan, p.AreaTujuan,
		tglBerangkat, tglKembali,
		p.MaksudTujuan,
		p.EstimasiBiaya,
	)

	waURL := "https://api.whatsapp.com/send?text=" + url.QueryEscape(waText)

	// Update status to WA_SENT_DIREKSI
	p.StatusFlow = StatusWASentDireksi
	p.UpdatedAt = time.Now()
	if err := config.DB.Save(p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui status: " + err.Error()})
	}
	enrichPengajuan(p)

	return c.JSON(fiber.Map{
		"message":      "Teks WA berhasil disiapkan, status diperbarui ke WA_SENT_DIREKSI",
		"text":         waText,
		"whatsapp_url": waURL,
		"wa_url":       waURL,
		"data":         p,
	})
}

// ============================================================
// PUT /api/pengajuan/:id/direksi-confirm
// HRGA menginput konfirmasi OK dari Direksi (WA_SENT_DIREKSI → CONFIRMED_DIREKSI)
// Body: { "konfirmasi_note": "OK / Disetujui" }
// ============================================================
func DireksiConfirm(c *fiber.Ctx) error {
	id := extractPengajuanID(c)
	p, err := getPengajuanWithKaryawan(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Pengajuan tidak ditemukan"})
	}
	if p.StatusFlow != StatusWASentDireksi && p.StatusFlow != StatusControlledHRD {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": fmt.Sprintf("Status harus '%s' untuk konfirmasi Direksi, saat ini: '%s'", StatusWASentDireksi, p.StatusFlow),
		})
	}

	type Input struct {
		KonfirmasiNote string `json:"konfirmasi_note"`
	}
	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Format input tidak valid"})
	}
	if input.KonfirmasiNote == "" {
		input.KonfirmasiNote = "OK / Disetujui"
	}

	now := time.Now()
	p.KonfirmasiDireksiNote = input.KonfirmasiNote
	p.TanggalKonfirmasiDireksi = &now
	p.StatusFlow = StatusConfirmedDireksi
	p.UpdatedAt = now

	if err := config.DB.Save(p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui status: " + err.Error()})
	}
	enrichPengajuan(p)
	return c.JSON(fiber.Map{"message": "Konfirmasi Direksi berhasil diinput", "data": p})
}

// ============================================================
// POST /api/pengajuan/:id/surat-tugas
// HRGA menerbitkan Surat Tugas resmi (CONFIRMED_DIREKSI → SURAT_TUGAS_ISSUED)
// Auto-generates NomorSuratTugas: ST/YYYY/MM/XXXX
// Returns full surat tugas payload for print view.
// ============================================================
func IssueSuratTugas(c *fiber.Ctx) error {
	id := extractPengajuanID(c)
	p, err := getPengajuanWithKaryawan(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Pengajuan tidak ditemukan"})
	}
	if p.StatusFlow != StatusConfirmedDireksi && p.StatusFlow != StatusWASentDireksi {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": fmt.Sprintf("Status harus '%s' untuk menerbitkan Surat Tugas, saat ini: '%s'", StatusConfirmedDireksi, p.StatusFlow),
		})
	}

	nomorSurat, err := generateNomorSuratTugas()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal generate nomor surat tugas"})
	}

	p.NomorSuratTugas = nomorSurat
	p.StatusFlow = StatusSuratTugasIssued
	p.UpdatedAt = time.Now()

	if err := config.DB.Save(p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan surat tugas: " + err.Error()})
	}
	enrichPengajuan(p)

	// Build printable surat tugas payload
	nama := p.NomorID
	dept := "-"
	jabatan := "-"
	if p.Karyawan != nil {
		nama = p.Karyawan.Nama
		dept = p.Karyawan.KodeDepartemen
		jabatan = p.Karyawan.KodeJabatan
	}

	suratTugasPayload := fiber.Map{
		"nomor_surat_tugas": nomorSurat,
		"no_pengajuan":      p.NoPengajuan,
		"nama":              nama,
		"nomor_id":          p.NomorID,
		"departemen":        dept,
		"jabatan":           jabatan,
		"area_tujuan":       p.AreaTujuan,
		"kota_tujuan":       p.KotaTujuan,
		"maksud_tujuan":     p.MaksudTujuan,
		"tanggal_berangkat": p.TanggalBerangkat.Format("02 January 2006"),
		"tanggal_kembali":   p.TanggalKembali.Format("02 January 2006"),
		"estimasi_biaya":    p.EstimasiBiaya,
		"catatan_hrd":       p.CatatanHRD,
		"tanggal_terbit":    time.Now().Format("02 January 2006"),
		"status_flow":       p.StatusFlow,
	}

	return c.JSON(fiber.Map{
		"message":     "Surat Tugas berhasil diterbitkan",
		"surat_tugas": suratTugasPayload,
		"data":        p,
	})
}