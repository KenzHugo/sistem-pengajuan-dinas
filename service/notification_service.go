package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sipadin/config"
	"sipadin/models"
	"sipadin/repository"
)

// NotificationService menangani seluruh alur notifikasi berjenjang SPPD dan rekan rombongan
type NotificationService interface {
	TriggerSPPDNotification(sppdID int, statusTahap string, actorID int, catatan string)
	FormatSPPDMessage(sppd *models.SppdPengajuan, recipientRole string, recipientName string, statusTahap string, actorName string, catatan string) string
	SendWhatsAppMessage(phone string, message string, recipientName string) error
	GenerateClickToChatURL(phone string, message string) string
}

type notificationService struct {
	sppdRepo repository.SppdRepository
	userRepo repository.UserRepository
}

// NewNotificationService membuat instance baru NotificationService
func NewNotificationService(sppdRepo repository.SppdRepository, userRepo repository.UserRepository) NotificationService {
	return &notificationService{
		sppdRepo: sppdRepo,
		userRepo: userRepo,
	}
}

// CleanPhoneNumber menstandarkan format nomor HP ke 628xxx
func CleanPhoneNumber(phone string) string {
	clean := strings.TrimSpace(phone)
	clean = strings.TrimPrefix(clean, "+")
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, "-", "")
	clean = strings.ReplaceAll(clean, "(", "")
	clean = strings.ReplaceAll(clean, ")", "")

	if strings.HasPrefix(clean, "08") {
		clean = "628" + clean[2:]
	} else if strings.HasPrefix(clean, "8") {
		clean = "628" + clean[1:]
	}
	return clean
}

// FormatSPPDMessage menyusun template pesan WhatsApp profesional dan terstruktur
func (n *notificationService) FormatSPPDMessage(
	sppd *models.SppdPengajuan,
	recipientRole string,
	recipientName string,
	statusTahap string,
	actorName string,
	catatan string,
) string {
	tglBerangkat := sppd.TglBerangkat.Format("02 Jan 2006")
	tglKembali := sppd.TglKembali.Format("02 Jan 2006")
	if sppd.TglBerangkat.IsZero() {
		tglBerangkat = "-"
	}
	if sppd.TglKembali.IsZero() {
		tglKembali = "-"
	}

	pemohonNama := "Karyawan"
	pemohonDept := "-"
	if sppd.Pemohon != nil {
		pemohonNama = sppd.Pemohon.NamaLengkap
		if sppd.Pemohon.KodeDepartemen != "" {
			pemohonDept = sppd.Pemohon.KodeDepartemen
		}
	}

	// Format daftar rekan kerja / rombongan
	rombonganLines := ""
	if len(sppd.Anggota) > 0 {
		var listRekan []string
		for idx, a := range sppd.Anggota {
			nama := "Rekan Kerja"
			nomorID := ""
			if a.User != nil {
				nama = a.User.NamaLengkap
				if a.User.NomorID != "" {
					nomorID = fmt.Sprintf(" (%s)", a.User.NomorID)
				}
			}
			listRekan = append(listRekan, fmt.Sprintf("   %d. %s%s", idx+1, nama, nomorID))
		}
		if len(listRekan) > 0 {
			rombonganLines = fmt.Sprintf("👥 *Rekan Kerja / Rombongan:*\n%s\n", strings.Join(listRekan, "\n"))
		}
	}

	// Status & Aksi Deskripsi
	var headerTitle, statusBadge, instruksiTambahan string
	switch statusTahap {
	case "SUBMITTED", "pending_manager":
		headerTitle = "🚀 *PENGAJUAN SPPD BARU*"
		statusBadge = "⏳ Menunggu Persetujuan Manager"
		if recipientRole == "Manager" {
			instruksiTambahan = "👉 *Aksi Diperlukan:* Mohon login ke sistem SiPaDin untuk memeriksa dan menyetujui pengajuan perjalanan dinas ini."
		} else if recipientRole == "Anggota" {
			instruksiTambahan = fmt.Sprintf("ℹ️ *Pemberitahuan:* Anda telah didaftarkan oleh *%s* sebagai anggota rombongan perjalanan dinas ini.", pemohonNama)
		} else {
			instruksiTambahan = "ℹ️ Pengajuan Anda telah diteruskan ke Manager untuk verifikasi tahap pertama."
		}

	case "MANAGER_APPROVED", "pending_hrga":
		headerTitle = "✅ *SPPD DISETUJUI MANAGER*"
		statusBadge = "⏳ Menunggu Verifikasi Fasilitas (HRGA)"
		if recipientRole == "HRGA" {
			instruksiTambahan = "👉 *Aksi Diperlukan:* Mohon verifikasi ketersediaan armada transportasi & akomodasi di sistem SiPaDin."
		} else if recipientRole == "Anggota" {
			instruksiTambahan = "ℹ️ SPPD telah disetujui Manager dan sedang dalam proses verifikasi fasilitas oleh HRGA."
		} else {
			instruksiTambahan = "ℹ️ Pengajuan Anda telah disetujui Manager dan diteruskan ke Tim HRGA."
		}

	case "HRGA_APPROVED", "pending_finance":
		headerTitle = "✅ *FASILITAS SPPD DIVERIFIKASI HRGA*"
		statusBadge = "⏳ Menunggu Persetujuan Anggaran (Finance)"
		if recipientRole == "Finance" {
			instruksiTambahan = "👉 *Aksi Diperlukan:* Mohon periksa dan setujui estimasi anggaran biaya di sistem SiPaDin."
		} else if recipientRole == "Anggota" {
			instruksiTambahan = "ℹ️ Fasilitas SPPD telah diverifikasi HRGA dan sedang menunggu persetujuan anggaran Finance."
		} else {
			instruksiTambahan = "ℹ️ Fasilitas SPPD telah diverifikasi HRGA dan diteruskan ke Tim Finance."
		}

	case "FINANCE_APPROVED", "ready_to_send":
		headerTitle = "✅ *ANGGARAN SPPD DISETUJUI FINANCE*"
		statusBadge = "📄 Siap Dikirim ke WhatsApp Direksi"
		if recipientRole == "HRGA" {
			instruksiTambahan = "👉 *Aksi Diperlukan:* Dokumen siap diajukan ke Direksi. Silakan kirim dokumen ke WhatsApp Direksi."
		} else if recipientRole == "Anggota" {
			instruksiTambahan = "ℹ️ Anggaran dinas telah disetujui Finance. Menunggu pengesahan akhir Direksi."
		} else {
			instruksiTambahan = "ℹ️ Anggaran dinas telah disetujui Finance. Dokumen segera diajukan ke Direksi."
		}

	case "SENT_TO_DIRECTOR", "approved":
		headerTitle = "🎉 *SPPD RESMI DISETUJUI DIREKSI*"
		statusBadge = "🌟 Disetujui Penuh (Approved)"
		if recipientRole == "Anggota" {
			instruksiTambahan = "🎉 *Selamat:* Perjalanan dinas rombongan Anda telah resmi disetujui Direksi. Harap persiapkan keperluan dinas sesuai jadwal."
		} else {
			instruksiTambahan = "🎉 *Selamat:* Perjalanan dinas Anda telah resmi disetujui Direksi. Silakan berkoordinasi dengan Finance & HRGA untuk keberangkatan."
		}

	case "REJECTED", "rejected":
		headerTitle = "❌ *PENGAJUAN SPPD DITOLAK*"
		statusBadge = "⛔ Ditolak"
		instruksiTambahan = "⚠️ Pengajuan perjalanan dinas ini tidak disetujui. Silakan cek catatan di bawah atau hubungi atasan terkait."

	default:
		headerTitle = "🔔 *NOTIFIKASI SISTEM SPPD*"
		statusBadge = statusTahap
		instruksiTambahan = "ℹ️ Status pengajuan SPPD telah diperbarui."
	}

	catatanLine := ""
	if strings.TrimSpace(catatan) != "" {
		catatanLine = fmt.Sprintf("📝 *Catatan:* %s\n", catatan)
	}

	tindakanOlehLine := ""
	if strings.TrimSpace(actorName) != "" {
		tindakanOlehLine = fmt.Sprintf("✍️ *Diproses Oleh:* %s\n", actorName)
	}

	return fmt.Sprintf(
		"%s\n"+
			"*PT ADIPRIMA SURAPRINTA*\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━\n\n"+
			"Halo *%s*,\n\n"+
			"📄 *No. SPPD:* %s\n"+
			"👤 *Pemohon Utama:* %s (Dept: %s)\n"+
			"%s"+
			"📍 *Tujuan:* %s - %s\n"+
			"🗓️ *Tanggal:* %s s/d %s (%d hari)\n"+
			"🚗 *Transportasi:* %s\n"+
			"🏨 *Akomodasi:* %s\n"+
			"💰 *Estimasi Biaya:* Rp %s\n\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━\n"+
			"📌 *Status Saat Ini:* %s\n"+
			"%s"+
			"%s\n"+
			"%s\n\n"+
			"🌐 *Sistem SPPD:* http://localhost:%s\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━\n"+
			"_Pesan otomatis dari Sistem Informasi Pengajuan Dinas (SiPaDin)_",
		headerTitle,
		recipientName,
		sppd.NoSppd,
		pemohonNama,
		pemohonDept,
		rombonganLines,
		sppd.KotaTujuan,
		sppd.InstansiTujuan,
		tglBerangkat,
		tglKembali,
		sppd.DurasiHari,
		sppd.Transportasi,
		sppd.Akomodasi,
		formatRupiah(sppd.EstimasiBiaya),
		statusBadge,
		tindakanOlehLine,
		catatanLine,
		instruksiTambahan,
		config.AppConfig.AppPort,
	)
}

// GenerateClickToChatURL membuat tautan https://wa.me/xxx?text=...
func (n *notificationService) GenerateClickToChatURL(phone string, message string) string {
	cleanPhone := CleanPhoneNumber(phone)
	return fmt.Sprintf("https://wa.me/%s?text=%s", cleanPhone, url.QueryEscape(message))
}

// SendWhatsAppMessage mengirim pesan melalui WA Gateway API atau mencatat log terstruktur
func (n *notificationService) SendWhatsAppMessage(phone string, message string, recipientName string) error {
	cleanPhone := CleanPhoneNumber(phone)
	if cleanPhone == "" {
		log.Printf("[Notification] [SKIP] Nomor telepon kosong untuk %s\n", recipientName)
		return nil
	}

	log.Printf("[Notification] 📨 [WA -> %s (%s)]:\n%s\n", cleanPhone, recipientName, message)

	cfg := config.AppConfig
	token := cfg.WAToken
	gatewayURL := cfg.WAGatewayURL

	// 1. Integrasi Gateway Fonnte jika token tersedia
	if token != "" && (gatewayURL == "" || strings.Contains(gatewayURL, "fonnte")) {
		targetURL := "https://api.fonnte.com/send"
		if gatewayURL != "" {
			targetURL = gatewayURL
		}

		payload := map[string]interface{}{
			"target":  cleanPhone,
			"message": message,
			"countryCode": "62",
		}
		jsonBytes, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(jsonBytes))
		if err == nil {
			req.Header.Set("Authorization", token)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: 5 * time.Second}
			resp, errResp := client.Do(req)
			if errResp == nil {
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				log.Printf("[Notification] [Fonnte Response] status: %d, body: %s\n", resp.StatusCode, string(body))
				return nil
			}
			log.Printf("[Notification] [Fonnte Error] Gagal kirim WA ke %s: %v\n", cleanPhone, errResp)
		}
	} else if gatewayURL != "" {
		// 2. Integrasi Generic Webhook Gateway
		payload := map[string]interface{}{
			"phone":     cleanPhone,
			"message":   message,
			"recipient": recipientName,
		}
		jsonBytes, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", gatewayURL, bytes.NewBuffer(jsonBytes))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			if cfg.WAApiKey != "" {
				req.Header.Set("X-API-KEY", cfg.WAApiKey)
			}
			client := &http.Client{Timeout: 5 * time.Second}
			resp, errResp := client.Do(req)
			if errResp == nil {
				defer resp.Body.Close()
				return nil
			}
			log.Printf("[Notification] [Webhook Error] Gagal kirim ke %s: %v\n", gatewayURL, errResp)
		}
	}

	return nil
}

// TriggerSPPDNotification memicu pengiriman notifikasi berjenjang ke Approver, Pemohon, dan Anggota Dinas
func (n *notificationService) TriggerSPPDNotification(sppdID int, statusTahap string, actorID int, catatan string) {
	// Jalankan secara asynchronous (goroutine) agar tidak memperlambat respon request HTTP
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Notification Panic] Terjadi kesalahan saat mengirim notifikasi SPPD #%d: %v\n", sppdID, r)
			}
		}()

		sppd, err := n.sppdRepo.FindByID(sppdID)
		if err != nil || sppd == nil {
			log.Printf("[Notification] SPPD #%d tidak ditemukan untuk notifikasi\n", sppdID)
			return
		}

		// Cari nama Actor / Approver yang memproses tindakan
		actorName := "Sistem"
		if actorID > 0 {
			if actor, errActor := n.userRepo.FindByID(actorID); errActor == nil && actor != nil {
				actorName = actor.NamaLengkap
			}
		}

		// ── 1. Notifikasi ke Pemohon Utama ──────────────────────────────────────
		if sppd.Pemohon != nil {
			pemohonPhone := sppd.Pemohon.NoWhatsapp
			if pemohonPhone == "" {
				pemohonPhone = "62895377328880" // fallback nomor tes
			}
			msgPemohon := n.FormatSPPDMessage(sppd, "Pemohon", sppd.Pemohon.NamaLengkap, statusTahap, actorName, catatan)
			_ = n.SendWhatsAppMessage(pemohonPhone, msgPemohon, sppd.Pemohon.NamaLengkap)
		}

		// ── 2. Notifikasi Tembusan ke Seluruh Rekan Kerja / Anggota Rombongan ──
		if len(sppd.Anggota) > 0 {
			for _, anggota := range sppd.Anggota {
				if anggota.User != nil {
					anggotaPhone := anggota.User.NoWhatsapp
					if anggotaPhone == "" {
						anggotaPhone = "62895377328880" // fallback nomor tes
					}
					msgAnggota := n.FormatSPPDMessage(sppd, "Anggota", anggota.User.NamaLengkap, statusTahap, actorName, catatan)
					_ = n.SendWhatsAppMessage(anggotaPhone, msgAnggota, anggota.User.NamaLengkap)
				}
			}
		}

		// ── 3. Notifikasi ke Approver Tahap Berikutnya ─────────────────────────
		var targetRole string
		switch statusTahap {
		case "SUBMITTED", "pending_manager":
			targetRole = "Manager"
		case "MANAGER_APPROVED", "pending_hrga":
			targetRole = "HRGA"
		case "HRGA_APPROVED", "pending_finance":
			targetRole = "Finance"
		case "FINANCE_APPROVED", "ready_to_send":
			targetRole = "HRGA" // HRGA menerima notif bahwa SPPD siap diteruskan ke Direksi
		}

		if targetRole != "" {
			approvers, errUsers := n.userRepo.FindByRole(targetRole)
			if errUsers == nil && len(approvers) > 0 {
				for _, appUser := range approvers {
					appPhone := appUser.NoWhatsapp
					if appPhone == "" {
						appPhone = "62895377328880"
					}
					msgApprover := n.FormatSPPDMessage(sppd, targetRole, appUser.NamaLengkap, statusTahap, actorName, catatan)
					_ = n.SendWhatsAppMessage(appPhone, msgApprover, fmt.Sprintf("%s (%s)", appUser.NamaLengkap, targetRole))
				}
			}
		}
	}()
}
