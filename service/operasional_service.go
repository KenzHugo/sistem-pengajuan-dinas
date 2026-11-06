package service

import (
	"fmt"
	"log"
	"sipadin/dto"
	"sipadin/repository"
	"strings"
	"time"
)

// OperasionalService mendefinisikan business logic untuk pengecekan jam operasional sistem
type OperasionalService interface {
	GetOperationalStatus(clientTimeStr string) (*dto.OperationalStatusResponse, error)
	ValidateAccess(identifier string, clientTimeStr string) (allowed bool, message string, isSpecial bool)
}

type operasionalService struct {
	repo repository.OperasionalRepository
}

// NewOperasionalService membuat instance baru OperasionalService
func NewOperasionalService(repo repository.OperasionalRepository) OperasionalService {
	return &operasionalService{repo: repo}
}

// Zona waktu WIB (Asia/Jakarta, UTC+7)
var locWIB = time.FixedZone("WIB", 7*3600)

var daysIndoMap = map[time.Weekday]string{
	time.Sunday:    "Minggu",
	time.Monday:    "Senin",
	time.Tuesday:   "Selasa",
	time.Wednesday: "Rabu",
	time.Thursday:  "Kamis",
	time.Friday:    "Jumat",
	time.Saturday:  "Sabtu",
}

// parseClientTime mengurai string waktu dari device pengguna dan mengonversinya ke zona WIB
func parseClientTime(clientTimeStr string, fallback time.Time) time.Time {
	clean := strings.TrimSpace(clientTimeStr)
	if clean == "" {
		return fallback
	}

	// Tangani kasus di mana karakter '+' pada URL query terkonversi menjadi spasi ' '
	if len(clean) >= 25 && clean[19] == ' ' && clean[22] == ':' {
		clean = clean[:19] + "+" + clean[20:]
	}

	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05+07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05 +0700",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"15:04:05",
		"15:04",
	}

	for _, f := range formats {
		if t, err := time.Parse(f, clean); err == nil {
			if f == "15:04" || f == "15:04:05" {
				return time.Date(fallback.Year(), fallback.Month(), fallback.Day(), t.Hour(), t.Minute(), t.Second(), 0, locWIB)
			}
			return t.In(locWIB)
		}
	}

	return fallback
}

// ValidateAccess memvalidasi secara cross-check waktu server (WIB) dan waktu device (WIB) terhadap tabel jamOperasional
func (s *operasionalService) ValidateAccess(identifier string, clientTimeStr string) (bool, string, bool) {
	// 1. Cek Spesial Akses terlebih dahulu
	if identifier != "" && s.repo.IsSpesialAkses(identifier) {
		log.Printf("[Operasional] Identifier '%s' memiliki SpesialAkses. Pembatasan jam kerja dilewati.", identifier)
		return true, "Akses diizinkan (Spesial Akses)", true
	}

	serverNow := time.Now().In(locWIB)
	deviceNow := parseClientTime(clientTimeStr, serverNow)

	// 2. Cross-check Hari Minggu (baik jika di device maupun di server adalah hari Minggu)
	if serverNow.Weekday() == time.Sunday || deviceNow.Weekday() == time.Sunday {
		return false, "bukan jam kerja, selamat istirahat 😊", false
	}

	hariIni := daysIndoMap[serverNow.Weekday()]

	// 3. Ambil data jamOperasional dari database
	item, err := s.repo.GetByHari(hariIni)
	if err != nil {
		log.Printf("[Operasional] Peringatan: Data jam operasional untuk hari %s tidak ditemukan: %v", hariIni, err)
		return true, "Dalam jam kerja", false
	}

	// 4. Cek jika status_libur true
	isLibur := (item.StatusLibur != nil && *item.StatusLibur) || strings.EqualFold(strings.TrimSpace(item.Hari), "Minggu")
	if isLibur {
		return false, "bukan jam kerja, selamat istirahat 😊", false
	}

	// 5. Cek kelengkapan kolom jam
	if item.JamMasuk == nil || item.JamSelesai == nil {
		return false, "bukan jam kerja, selamat istirahat 😊", false
	}

	masukSecs := item.JamMasuk.Hour()*3600 + item.JamMasuk.Minute()*60 + item.JamMasuk.Second()
	selesaiSecs := item.JamSelesai.Hour()*3600 + item.JamSelesai.Minute()*60 + item.JamSelesai.Second()

	// A. Validasi Jam di Device Pengguna (WIB)
	devSecs := deviceNow.Hour()*3600 + deviceNow.Minute()*60 + deviceNow.Second()
	if devSecs < masukSecs {
		log.Printf("[Operasional] Jam device (%s WIB) belum masuk jam kerja", deviceNow.Format("15:04"))
		return false, "belum masuk jam kerja", false
	}
	if devSecs > selesaiSecs {
		log.Printf("[Operasional] Jam device (%s WIB) sudah melewati jam kerja", deviceNow.Format("15:04"))
		return false, "sudah melewati jam kerja", false
	}

	// B. Validasi Jam di Server (WIB)
	srvSecs := serverNow.Hour()*3600 + serverNow.Minute()*60 + serverNow.Second()
	if srvSecs < masukSecs {
		log.Printf("[Operasional] Jam server (%s WIB) belum masuk jam kerja", serverNow.Format("15:04"))
		return false, "belum masuk jam kerja", false
	}
	if srvSecs > selesaiSecs {
		log.Printf("[Operasional] Jam server (%s WIB) sudah melewati jam kerja", serverNow.Format("15:04"))
		return false, "sudah melewati jam kerja", false
	}

	return true, "dalam jam kerja", false
}

// GetOperationalStatus mengembalikan informasi lengkap status operasional hasil cross-check server & device
func (s *operasionalService) GetOperationalStatus(clientTimeStr string) (*dto.OperationalStatusResponse, error) {
	serverNow := time.Now().In(locWIB)
	deviceNow := parseClientTime(clientTimeStr, serverNow)

	isSunday := serverNow.Weekday() == time.Sunday || deviceNow.Weekday() == time.Sunday
	hariIni := daysIndoMap[serverNow.Weekday()]

	item, err := s.repo.GetByHari(hariIni)
	jamMasukStr := ""
	jamSelesaiStr := ""
	statusLibur := isSunday

	if err == nil && item != nil {
		if item.StatusLibur != nil && *item.StatusLibur {
			statusLibur = true
		}
		if item.JamMasuk != nil {
			jamMasukStr = fmt.Sprintf("%02d:%02d", item.JamMasuk.Hour(), item.JamMasuk.Minute())
		}
		if item.JamSelesai != nil {
			jamSelesaiStr = fmt.Sprintf("%02d:%02d", item.JamSelesai.Hour(), item.JamSelesai.Minute())
		}
	}

	allowed, msg, _ := s.ValidateAccess("", clientTimeStr)

	resp := &dto.OperationalStatusResponse{
		IsOpen:      allowed,
		Message:     msg,
		Hari:        hariIni,
		JamMasuk:    jamMasukStr,
		JamSelesai:  jamSelesaiStr,
		StatusLibur: statusLibur,
		ServerTime:  serverNow.Format("15:04"),
		DeviceTime:  deviceNow.Format("15:04"),
		CurrentTime: serverNow.Format("15:04"),
		IsSunday:    isSunday,
	}

	return resp, nil
}
