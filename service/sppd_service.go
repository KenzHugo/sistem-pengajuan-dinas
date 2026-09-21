package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"sipadin/config"
	"sipadin/dto"
	"sipadin/models"
	"sipadin/repository"
)

// SppdService mendefinisikan business logic pengajuan SPPD
type SppdService interface {
	Create(idPemohon int, req dto.CreateSppdRequest) (*models.SppdPengajuan, error)
	Submit(idSppd, idPemohon int) (*models.SppdPengajuan, error)
	GetByID(id int) (*models.SppdPengajuan, error)
	GetMyList(idPemohon int, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error)
	GetAllList(filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error)
	GetByStatus(statuses []string, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error)
	Update(idSppd, idPemohon int, req dto.UpdateSppdRequest) (*models.SppdPengajuan, error)
	Delete(idSppd, idPemohon int) error
	HitungEstimasiBiaya(idPemohon int, req dto.KalkulasiBiayaRequest) (*dto.KalkulasiBiayaResponse, error)
	HitungRincianBiayaSPPD(nomorID string, jarak float64, jamBerangkat, jamKembali, tglBerangkat, tglKembali, kotaTujuan, transportasi, akomodasi string, rekanKerja []string) (*dto.KalkulasiBiayaResponse, error)
}

type sppdService struct {
	sppdRepo    repository.SppdRepository
	anggotaRepo repository.AnggotaRepository
	tarifRepo   repository.TarifRepository
	userRepo    repository.UserRepository
	notifSvc    NotificationService
}

func NewSppdService(
	sppdRepo repository.SppdRepository,
	anggotaRepo repository.AnggotaRepository,
	tarifRepo repository.TarifRepository,
	userRepo repository.UserRepository,
	notifSvc NotificationService,
) SppdService {
	return &sppdService{
		sppdRepo:    sppdRepo,
		anggotaRepo: anggotaRepo,
		tarifRepo:   tarifRepo,
		userRepo:    userRepo,
		notifSvc:    notifSvc,
	}
}

// resolveUserID mencari atau menyinkronkan user ID berdasarkan ID integer, Nomor ID, atau Nama (termasuk format datalist)
func (s *sppdService) resolveUserID(identifier string) int {
	clean := strings.TrimSpace(identifier)
	if clean == "" {
		return 0
	}

	extractedID := clean
	extractedNama := clean

	// Jika format berupa "Nama (NomorID) - Dept: XXX", ekstrak dengan LastIndex agar aman untuk nama dengan tanda kurung seperti "KHOIRI (B) (01001200208)"
	if lastOpen := strings.LastIndex(clean, "("); lastOpen != -1 {
		if lastClose := strings.LastIndex(clean, ")"); lastClose != -1 && lastClose > lastOpen {
			extractedID = strings.TrimSpace(clean[lastOpen+1 : lastClose])
			extractedNama = strings.TrimSpace(clean[:lastOpen])
		}
	}

	cleanExpr := "TRIM(REPLACE(REPLACE(nomor_id, '\r', ''), '\n', ''))"
	cleanNama := "TRIM(REPLACE(REPLACE(nama_lengkap, '\r', ''), '\n', ''))"

	// 1. Coba cari langsung di tabel users
	var user models.User
	err := config.DB.Where(
		cleanExpr+" = ? OR "+cleanExpr+" = ? OR LOWER("+cleanNama+") = LOWER(?) OR LOWER("+cleanNama+") = LOWER(?) OR "+cleanNama+" LIKE ?",
		extractedID, clean, extractedNama, clean, "%"+extractedNama+"%",
	).First(&user).Error
	if err == nil && user.IDUser > 0 {
		return user.IDUser
	}

	// 2. Coba cari di tabel karyawan_copy1 dan sinkronkan ke users jika belum ada
	var karyawan models.Karyawan
	kCleanID := "TRIM(REPLACE(REPLACE(NomorID, '\r', ''), '\n', ''))"
	kCleanNama := "TRIM(REPLACE(REPLACE(Nama, '\r', ''), '\n', ''))"
	errK := config.DB.Table("karyawan_copy1").Where(
		kCleanID+" = ? OR "+kCleanID+" = ? OR LOWER("+kCleanNama+") = LOWER(?) OR LOWER("+kCleanNama+") = LOWER(?) OR "+kCleanNama+" LIKE ?",
		extractedID, clean, extractedNama, clean, "%"+extractedNama+"%",
	).First(&karyawan).Error
	if errK == nil {
		cID := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(karyawan.NomorID, "\r", ""), "\n", ""))
		cNama := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(karyawan.Nama, "\r", ""), "\n", ""))

		// Cek lagi apakah user dengan nomor ID / nama ini sudah ada di users
		var existing models.User
		if errE := config.DB.Where(cleanExpr+" = ? OR LOWER("+cleanNama+") = LOWER(?)", cID, cNama).First(&existing).Error; errE == nil {
			return existing.IDUser
		}

		newUser := models.User{
			NomorID:        cID,
			NamaLengkap:    cNama,
			KodeDepartemen: karyawan.KodeDepartemen,
			KodeJabatan:    karyawan.KodeJabatan,
			RoleSistem:     "Karyawan",
			PasswordHash:   karyawan.Password,
			StatusAktif:    1,
		}
		if errCreate := config.DB.Create(&newUser).Error; errCreate == nil {
			return newUser.IDUser
		}
	}

	// 3. Fallback jika merupakan integer ID user murni
	if id, err := strconv.Atoi(extractedID); err == nil && id > 0 {
		var u models.User
		if err := config.DB.First(&u, id).Error; err == nil {
			return u.IDUser
		}
	}

	return 0
}

// Create membuat pengajuan SPPD baru dengan status draft
func (s *sppdService) Create(idPemohon int, req dto.CreateSppdRequest) (*models.SppdPengajuan, error) {
	tglBerangkat, err := time.Parse("2006-01-02", req.TglBerangkat)
	if err != nil {
		return nil, errors.New("format tgl_berangkat tidak valid, gunakan YYYY-MM-DD")
	}
	tglKembali, err := time.Parse("2006-01-02", req.TglKembali)
	if err != nil {
		return nil, errors.New("format tgl_kembali tidak valid, gunakan YYYY-MM-DD")
	}
	if tglKembali.Before(tglBerangkat) {
		return nil, errors.New("tgl_kembali tidak boleh sebelum tgl_berangkat")
	}

	durasi := int(tglKembali.Sub(tglBerangkat).Hours()/24) + 1

	jamBerangkat := strings.TrimSpace(req.JamBerangkat)
	if jamBerangkat == "" {
		jamBerangkat = "08:00"
	}
	jamKembali := strings.TrimSpace(req.JamKembali)
	if jamKembali == "" {
		jamKembali = "17:00"
	}

	calc, _ := s.HitungEstimasiBiaya(idPemohon, dto.KalkulasiBiayaRequest{
		KotaTujuan:   req.KotaTujuan,
		TglBerangkat: req.TglBerangkat,
		TglKembali:   req.TglKembali,
		JamBerangkat: jamBerangkat,
		JamKembali:   jamKembali,
		JarakKM:      req.JarakKM,
		Transportasi: req.Transportasi,
		Akomodasi:    req.Akomodasi,
		RekanKerja:   req.RekanKerja,
	})

	// Backend WAJIB menghitung mandiri secara objektif (mengabaikan input manual karyawan)
	biayaMakan := 0.0
	biayaTransport := 0.0
	biayaUangSaku := 0.0
	biayaTiket := 0.0
	biayaHotel := 0.0
	estimasiBiaya := 0.0
	durasiJam := 0.0
	plafonMakan := 0.0
	plafonTransport := 0.0
	plafonUangSaku := 0.0
	plafonTiket := 0.0
	plafonHotel := 0.0
	statusPlafon := "SESUAI"

	if calc != nil {
		durasiJam = calc.DurasiJam
		if calc.DurasiHari > 0 {
			durasi = calc.DurasiHari
		}
		plafonMakan = calc.PlafonMakan
		plafonTransport = calc.PlafonTransport
		plafonUangSaku = calc.PlafonUangSaku
		plafonTiket = calc.PlafonTiket
		plafonHotel = calc.PlafonHotel

		biayaMakan = calc.PlafonMakan
		biayaTransport = calc.PlafonTransport
		biayaUangSaku = calc.PlafonUangSaku
		biayaHotel = calc.PlafonHotel
		biayaTiket = calc.PlafonTiket
		estimasiBiaya = calc.EstimasiBiayaTotal
	}

	rincianBiaya := req.RincianBiaya
	if rincianBiaya == "" && calc != nil {
		rincianBiaya = fmt.Sprintf(
			"Kalkulasi Sistem: Makan %s | Trans %s | Saku %s | Tiket %s | Hotel %s | Total: %s (Status: %s)",
			formatRupiahGo(biayaMakan), formatRupiahGo(biayaTransport), formatRupiahGo(biayaUangSaku),
			formatRupiahGo(biayaTiket), formatRupiahGo(biayaHotel), formatRupiahGo(estimasiBiaya), statusPlafon,
		)
	}

	sppd := &models.SppdPengajuan{
		IDPemohon:       idPemohon,
		KotaTujuan:      req.KotaTujuan,
		InstansiTujuan:  req.InstansiTujuan,
		MaksudTujuan:    req.MaksudTujuan,
		TglBerangkat:    tglBerangkat,
		TglKembali:      tglKembali,
		JamBerangkat:    jamBerangkat,
		JamKembali:      jamKembali,
		DurasiHari:      durasi,
		DurasiJam:       durasiJam,
		JarakKM:         req.JarakKM,
		Transportasi:    req.Transportasi,
		Akomodasi:       req.Akomodasi,
		BiayaMakan:      biayaMakan,
		BiayaTransport:  biayaTransport,
		BiayaUangSaku:   biayaUangSaku,
		BiayaTiket:      biayaTiket,
		BiayaHotel:      biayaHotel,
		PlafonMakan:     plafonMakan,
		PlafonTransport: plafonTransport,
		PlafonUangSaku:  plafonUangSaku,
		PlafonTiket:     plafonTiket,
		PlafonHotel:     plafonHotel,
		StatusPlafon:    statusPlafon,
		EstimasiBiaya:   estimasiBiaya,
		RincianBiaya:    rincianBiaya,
		LampiranFile:    req.LampiranFile,
		StatusDokumen:   "draft",
	}

	if err := s.sppdRepo.Create(sppd); err != nil {
		return nil, errors.New("gagal menyimpan pengajuan")
	}

	// Simpan anggota / rekan kerja jika ada
	processedUsers := make(map[int]bool)
	if idPemohon > 0 {
		processedUsers[idPemohon] = true // Pemohon utama tidak perlu dobel di rombongan
	}

	for _, rekan := range req.RekanKerja {
		targetID := s.resolveUserID(rekan)
		if targetID > 0 && !processedUsers[targetID] {
			processedUsers[targetID] = true
			anggota := &models.SppdAnggota{
				IDSppd:            sppd.IDSppd,
				IDUser:            targetID,
				KriteriaFasilitas: "Karyawan",
				UangMakan:         0,
				UangSaku:          0,
			}
			_ = s.anggotaRepo.Create(anggota)
		}
	}

	for _, a := range req.Anggota {
		if a.IDUser > 0 && !processedUsers[a.IDUser] {
			processedUsers[a.IDUser] = true
			anggota := &models.SppdAnggota{
				IDSppd:            sppd.IDSppd,
				IDUser:            a.IDUser,
				KriteriaFasilitas: a.KriteriaFasilitas,
				UangMakan:         a.UangMakan,
				UangSaku:          a.UangSaku,
			}
			_ = s.anggotaRepo.Create(anggota)
		}
	}

	return s.sppdRepo.FindByID(sppd.IDSppd)
}

// Submit mengubah status SPPD dari draft menjadi pending_manager (mengajukan ke Manager)
func (s *sppdService) Submit(idSppd, idPemohon int) (*models.SppdPengajuan, error) {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if sppd.IDPemohon != idPemohon {
		return nil, errors.New("anda tidak memiliki akses ke pengajuan ini")
	}
	if sppd.StatusDokumen != "draft" {
		return nil, errors.New("hanya pengajuan berstatus draft yang bisa diajukan")
	}
	if sppd.NoSppd == "" {
		if noSppd, errGen := s.sppdRepo.GenerateNoSppd(); errGen == nil && noSppd != "" {
			sppd.NoSppd = noSppd
		}
	}

	sppd.StatusDokumen = "pending_manager"
	if err := s.sppdRepo.Update(sppd); err != nil {
		return nil, errors.New("gagal mengajukan SPPD")
	}

	// Kirim notifikasi berjenjang ke Manager & Anggota rombongan
	if s.notifSvc != nil {
		go s.notifSvc.TriggerSPPDNotification(idSppd, "SUBMITTED", idPemohon, "Pengajuan baru diajukan oleh pemohon")
	}

	return s.sppdRepo.FindByID(idSppd)
}

// GetByID mengambil detail SPPD beserta relasi
func (s *sppdService) GetByID(id int) (*models.SppdPengajuan, error) {
	sppd, err := s.sppdRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	return sppd, nil
}

// GetMyList mengambil daftar SPPD milik user yang sedang login
func (s *sppdService) GetMyList(idPemohon int, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error) {
	filter.IDPemohon = idPemohon
	return s.sppdRepo.FindByPemohon(idPemohon, filter)
}

// GetAllList mengambil semua SPPD (untuk HRGA / Finance / Direksi)
func (s *sppdService) GetAllList(filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error) {
	return s.sppdRepo.FindAll(filter)
}

// GetByStatus mengambil SPPD berdasarkan beberapa status sekaligus
func (s *sppdService) GetByStatus(statuses []string, filter dto.SppdListFilter) ([]models.SppdPengajuan, int64, error) {
	return s.sppdRepo.FindByStatus(statuses, filter)
}

// Update mengedit pengajuan yang masih berstatus draft
func (s *sppdService) Update(idSppd, idPemohon int, req dto.UpdateSppdRequest) (*models.SppdPengajuan, error) {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if sppd.IDPemohon != idPemohon {
		return nil, errors.New("anda tidak memiliki akses ke pengajuan ini")
	}
	if sppd.StatusDokumen != "draft" {
		return nil, errors.New("hanya pengajuan berstatus draft yang bisa diedit")
	}

	if req.KotaTujuan != "" {
		sppd.KotaTujuan = req.KotaTujuan
	}
	if req.InstansiTujuan != "" {
		sppd.InstansiTujuan = req.InstansiTujuan
	}
	if req.MaksudTujuan != "" {
		sppd.MaksudTujuan = req.MaksudTujuan
	}
	if req.TglBerangkat != "" {
		t, _ := time.Parse("2006-01-02", req.TglBerangkat)
		sppd.TglBerangkat = t
	}
	if req.TglKembali != "" {
		t, _ := time.Parse("2006-01-02", req.TglKembali)
		sppd.TglKembali = t
	}
	if req.JamBerangkat != "" {
		sppd.JamBerangkat = req.JamBerangkat
	}
	if req.JamKembali != "" {
		sppd.JamKembali = req.JamKembali
	}
	if req.JarakKM > 0 {
		sppd.JarakKM = req.JarakKM
	}
	if req.Transportasi != "" {
		sppd.Transportasi = req.Transportasi
	}
	if req.Akomodasi != "" {
		sppd.Akomodasi = req.Akomodasi
	}
	if req.LampiranFile != "" {
		sppd.LampiranFile = req.LampiranFile
	}

	sppd.DurasiHari = int(sppd.TglKembali.Sub(sppd.TglBerangkat).Hours()/24) + 1
	if sppd.DurasiHari < 1 {
		sppd.DurasiHari = 1
	}

	// Hitung ulang estimasi biaya dan batas plafon secara mandiri di server
	calc, _ := s.HitungEstimasiBiaya(idPemohon, dto.KalkulasiBiayaRequest{
		KotaTujuan:   sppd.KotaTujuan,
		TglBerangkat: sppd.TglBerangkat.Format("2006-01-02"),
		TglKembali:   sppd.TglKembali.Format("2006-01-02"),
		JamBerangkat: sppd.JamBerangkat,
		JamKembali:   sppd.JamKembali,
		JarakKM:      sppd.JarakKM,
		Transportasi: sppd.Transportasi,
		Akomodasi:    sppd.Akomodasi,
		RekanKerja:   req.RekanKerja,
	})

	if calc != nil {
		sppd.DurasiJam = calc.DurasiJam
		sppd.PlafonMakan = calc.PlafonMakan
		sppd.PlafonTransport = calc.PlafonTransport
		sppd.PlafonUangSaku = calc.PlafonUangSaku
		sppd.PlafonTiket = calc.PlafonTiket
		sppd.PlafonHotel = calc.PlafonHotel

		sppd.BiayaMakan = calc.PlafonMakan
		sppd.BiayaTransport = calc.PlafonTransport
		sppd.BiayaUangSaku = calc.PlafonUangSaku
		sppd.BiayaTiket = calc.PlafonTiket
		sppd.BiayaHotel = calc.PlafonHotel
		sppd.EstimasiBiaya = calc.EstimasiBiayaTotal
		sppd.StatusPlafon = "SESUAI"

		if req.RincianBiaya != "" {
			sppd.RincianBiaya = req.RincianBiaya
		} else {
			sppd.RincianBiaya = fmt.Sprintf(
				"Kalkulasi Sistem: Makan %s | Trans %s | Saku %s | Tiket %s | Hotel %s | Total: %s (Status: %s)",
				formatRupiahGo(sppd.BiayaMakan), formatRupiahGo(sppd.BiayaTransport), formatRupiahGo(sppd.BiayaUangSaku),
				formatRupiahGo(sppd.BiayaTiket), formatRupiahGo(sppd.BiayaHotel), formatRupiahGo(sppd.EstimasiBiaya), sppd.StatusPlafon,
			)
		}
	}

	if err := s.sppdRepo.Update(sppd); err != nil {
		return nil, errors.New("gagal mengupdate pengajuan")
	}

	// Sync ulang anggota / rekan kerja jika dikirim
	if len(req.RekanKerja) > 0 || len(req.Anggota) > 0 {
		_ = s.anggotaRepo.DeleteBySppd(idSppd)
		processedUsers := make(map[int]bool)
		if idPemohon > 0 {
			processedUsers[idPemohon] = true
		}

		for _, rekan := range req.RekanKerja {
			targetID := s.resolveUserID(rekan)
			if targetID > 0 && !processedUsers[targetID] {
				processedUsers[targetID] = true
				anggota := &models.SppdAnggota{
					IDSppd:            idSppd,
					IDUser:            targetID,
					KriteriaFasilitas: "Karyawan",
					UangMakan:         0,
					UangSaku:          0,
				}
				_ = s.anggotaRepo.Create(anggota)
			}
		}

		for _, a := range req.Anggota {
			if a.IDUser > 0 && !processedUsers[a.IDUser] {
				processedUsers[a.IDUser] = true
				anggota := &models.SppdAnggota{
					IDSppd:            idSppd,
					IDUser:            a.IDUser,
					KriteriaFasilitas: a.KriteriaFasilitas,
					UangMakan:         a.UangMakan,
					UangSaku:          a.UangSaku,
				}
				_ = s.anggotaRepo.Create(anggota)
			}
		}
	}

	return s.sppdRepo.FindByID(idSppd)
}

// Delete menghapus pengajuan yang masih berstatus draft
func (s *sppdService) Delete(idSppd, idPemohon int) error {
	sppd, err := s.sppdRepo.FindByID(idSppd)
	if err != nil {
		return errors.New("pengajuan tidak ditemukan")
	}
	if sppd.IDPemohon != idPemohon {
		return errors.New("anda tidak memiliki akses ke pengajuan ini")
	}
	if sppd.StatusDokumen != "draft" {
		return errors.New("hanya pengajuan berstatus draft yang bisa dihapus")
	}

	_ = s.anggotaRepo.DeleteBySppd(idSppd)
	return s.sppdRepo.Delete(idSppd)
}

// DetermineArea menentukan kategori Area tarif berdasarkan nama kota tujuan
func DetermineArea(kota string) string {
	k := strings.ToLower(strings.TrimSpace(kota))
	if k == "" {
		return "Jawa"
	}

	// 1. Jakarta / Luar Jawa (Prioritaskan cek kota / wilayah di Indonesia)
	if strings.Contains(k, "jakarta") || strings.Contains(k, "luar jawa") ||
		strings.Contains(k, "sumatera") || strings.Contains(k, "medan") || strings.Contains(k, "palembang") ||
		strings.Contains(k, "padang") || strings.Contains(k, "aceh") || strings.Contains(k, "pekanbaru") ||
		strings.Contains(k, "lampung") || strings.Contains(k, "jambi") || strings.Contains(k, "bengkulu") ||
		strings.Contains(k, "riau") || strings.Contains(k, "batam") || strings.Contains(k, "bangka") ||
		strings.Contains(k, "belitung") || strings.Contains(k, "pangkal pinang") ||
		strings.Contains(k, "kalimantan") || strings.Contains(k, "balikpapan") || strings.Contains(k, "samarinda") ||
		strings.Contains(k, "pontianak") || strings.Contains(k, "banjarmasin") || strings.Contains(k, "palangkaraya") ||
		strings.Contains(k, "tarakan") || strings.Contains(k, "sulawesi") || strings.Contains(k, "makassar") ||
		strings.Contains(k, "manado") || strings.Contains(k, "palu") || strings.Contains(k, "kendari") ||
		strings.Contains(k, "gorontalo") || strings.Contains(k, "mamuju") ||
		strings.Contains(k, "bali") || strings.Contains(k, "denpasar") || strings.Contains(k, "kuta") ||
		strings.Contains(k, "lombok") || strings.Contains(k, "mataram") || strings.Contains(k, "ntb") ||
		strings.Contains(k, "ntt") || strings.Contains(k, "kupang") || strings.Contains(k, "labuan bajo") ||
		strings.Contains(k, "maluku") || strings.Contains(k, "ambon") || strings.Contains(k, "ternate") ||
		strings.Contains(k, "papua") || strings.Contains(k, "jayapura") || strings.Contains(k, "manokwari") ||
		strings.Contains(k, "sorong") || strings.Contains(k, "merauke") || strings.Contains(k, "timika") {
		return "Jakarta / Luar Jawa"
	}

	// 2. Amerika / Singapura
	if strings.Contains(k, "singapura") || strings.Contains(k, "singapore") || strings.Contains(k, "amerika") ||
		strings.Contains(k, "united states") || strings.Contains(k, "new york") ||
		strings.Contains(k, "los angeles") || strings.Contains(k, "san francisco") || strings.Contains(k, "chicago") ||
		strings.Contains(k, "seattle") || strings.Contains(k, "boston") || strings.Contains(k, "california") ||
		strings.Contains(k, "texas") || strings.Contains(k, "washington") || k == "usa" || k == "us" || strings.HasPrefix(k, "usa ") || strings.HasSuffix(k, " usa") || strings.Contains(k, " usa ") {
		return "Amerika / Singapura"
	}

	// 3. Jepang
	if strings.Contains(k, "jepang") || strings.Contains(k, "japan") || strings.Contains(k, "tokyo") ||
		strings.Contains(k, "osaka") || strings.Contains(k, "kyoto") || strings.Contains(k, "nagoya") ||
		strings.Contains(k, "fukuoka") || strings.Contains(k, "sapporo") || strings.Contains(k, "yokohama") {
		return "Jepang"
	}

	// 4. Eropa
	if strings.Contains(k, "eropa") || strings.Contains(k, "europe") || strings.Contains(k, "jerman") ||
		strings.Contains(k, "germany") || strings.Contains(k, "belanda") || strings.Contains(k, "netherlands") ||
		strings.Contains(k, "prancis") || strings.Contains(k, "france") || strings.Contains(k, "paris") ||
		strings.Contains(k, "inggris") || strings.Contains(k, "london") || strings.Contains(k, "united kingdom") ||
		strings.Contains(k, "italia") || strings.Contains(k, "italy") || strings.Contains(k, "roma") ||
		strings.Contains(k, "spanyol") || strings.Contains(k, "spain") || strings.Contains(k, "swiss") ||
		strings.Contains(k, "switzerland") || strings.Contains(k, "austria") || strings.Contains(k, "belgia") {
		return "Eropa"
	}

	// 5. Asia
	if strings.Contains(k, "asia") || strings.Contains(k, "malaysia") || strings.Contains(k, "kuala lumpur") ||
		strings.Contains(k, "thailand") || strings.Contains(k, "bangkok") || strings.Contains(k, "vietnam") ||
		strings.Contains(k, "hanoi") || strings.Contains(k, "filipina") || strings.Contains(k, "philippines") ||
		strings.Contains(k, "manila") || strings.Contains(k, "china") || strings.Contains(k, "tiongkok") ||
		strings.Contains(k, "beijing") || strings.Contains(k, "shanghai") || strings.Contains(k, "korea") ||
		strings.Contains(k, "seoul") || strings.Contains(k, "taiwan") || strings.Contains(k, "taipei") ||
		strings.Contains(k, "hongkong") || strings.Contains(k, "hong kong") || strings.Contains(k, "india") ||
		strings.Contains(k, "delhi") || strings.Contains(k, "mumbai") || strings.Contains(k, "brunei") ||
		strings.Contains(k, "kamboja") || strings.Contains(k, "cambodia") || strings.Contains(k, "laos") ||
		strings.Contains(k, "myanmar") || strings.Contains(k, "australia") || strings.Contains(k, "sydney") {
		return "Asia"
	}

	// 6. Default: Jawa (Surabaya, Gresik, Sidoarjo, Mojokerto, Malang, Semarang, Solo, Yogyakarta, Bandung, Cirebon, etc.)
	return "Jawa"
}

// DetermineJabatan memetakan kode jabatan / role pemohon ke kategori master tarif
func DetermineJabatan(user *models.User) string {
	if user == nil {
		return "Operator - Karu"
	}

	role := strings.ToLower(strings.TrimSpace(user.RoleSistem))
	kode := strings.TrimSpace(user.KodeJabatan)
	nama := strings.ToLower(strings.TrimSpace(user.NamaLengkap))

	// Direktur - Dirut (Role Direksi, Kode >= 31)
	if role == "direksi" || strings.Contains(nama, "direktur") || strings.Contains(nama, "dirut") {
		return "Direktur - Dirut"
	}
	if kodeInt, err := strconv.Atoi(kode); err == nil && kodeInt >= 31 {
		return "Direktur - Dirut"
	}

	// GM - Wadir
	if strings.Contains(nama, "general manager") || strings.Contains(nama, "gm") || strings.Contains(nama, "wadir") {
		return "GM - Wadir"
	}
	if kode == "29" || kode == "30" {
		return "GM - Wadir"
	}

	// Assman - Manager
	if role == "manager" || strings.Contains(nama, "manager") || strings.Contains(nama, "assman") || strings.Contains(nama, "asisten manager") {
		return "Assman - Manager"
	}
	if kode == "27" || kode == "28" {
		return "Assman - Manager"
	}

	// Assup - Supervisor
	if strings.Contains(nama, "supervisor") || strings.Contains(nama, "spv") || strings.Contains(nama, "assup") || strings.Contains(nama, "asisten supervisor") {
		return "Assup - Supervisor"
	}
	if kode == "25" || kode == "26" {
		return "Assup - Supervisor"
	}

	// Default: Operator - Karu
	return "Operator - Karu"
}

// MapKodeJabatanToGroup memetakan kode jabatan numerik ke kategori GroupDinasLuar
func MapKodeJabatanToGroup(kode string) string {
	k := strings.TrimSpace(kode)
	if k == "" {
		return "Operator - Karu"
	}
	if n, err := strconv.Atoi(k); err == nil {
		if n >= 31 {
			return "Direktur - Dirut"
		} else if n >= 29 {
			return "GM - Wadir"
		} else if n >= 27 {
			return "Assman - Manager"
		} else if n >= 25 {
			return "Assup - Supervisor"
		}
	}
	return "Operator - Karu"
}

// HitungRincianBiayaSPPD adalah helper kalkulator biaya SPPD objektif berdasarkan jarak, jam, dan plafon jabatan
func (s *sppdService) HitungRincianBiayaSPPD(nomorID string, jarak float64, jamBerangkat, jamKembali, tglBerangkat, tglKembali, kotaTujuan, transportasi, akomodasi string, rekanKerja []string) (*dto.KalkulasiBiayaResponse, error) {
	req := dto.KalkulasiBiayaRequest{
		NomorID:      nomorID,
		KotaTujuan:   kotaTujuan,
		TglBerangkat: tglBerangkat,
		TglKembali:   tglKembali,
		JamBerangkat: jamBerangkat,
		JamKembali:   jamKembali,
		JarakKM:      jarak,
		Transportasi: transportasi,
		Akomodasi:    akomodasi,
		RekanKerja:   rekanKerja,
	}
	return s.HitungEstimasiBiaya(0, req)
}

// PersonilDinasTarif menampung data tarif individual hasil JOIN karyawan -> jabatan -> tarif
type PersonilDinasTarif struct {
	NomorID        string  `gorm:"column:NomorID"`
	Nama           string  `gorm:"column:Nama"`
	KodeJabatan    string  `gorm:"column:KodeJabatan"`
	NamaJabatan    string  `gorm:"column:NamaJabatan"`
	GroupDinasLuar string  `gorm:"column:GroupDinasLuar"`
	Area           string  `gorm:"column:Area"`
	MataUang       string  `gorm:"column:MataUang"`
	UangMakan      float64 `gorm:"column:UangMakan"`
	UangSaku       float64 `gorm:"column:UangSaku"`
	Hotel          float64 `gorm:"column:Hotel"`
	Transport      float64 `gorm:"column:Transport"`
	Tiket          float64 `gorm:"column:Tiket"`
}

// HitungEstimasiBiaya menghitung estimasi biaya SPPD secara dinamis per kepala (Pemohon + Seluruh Partner)
// Menggunakan relasi tabel: karyawan_copy1 -> jabatan (KodeJabatan) -> perjalanandinas_tarif (GroupDinasLuar)
func (s *sppdService) HitungEstimasiBiaya(idPemohon int, req dto.KalkulasiBiayaRequest) (*dto.KalkulasiBiayaResponse, error) {
	// 1. Normalisasi jam & tanggal perjalanan
	jamB := strings.TrimSpace(req.JamBerangkat)
	if jamB == "" {
		jamB = "08:00"
	}
	jamK := strings.TrimSpace(req.JamKembali)
	if jamK == "" {
		jamK = "17:00"
	}

	tglBStr := strings.TrimSpace(req.TglBerangkat)
	if tglBStr == "" {
		tglBStr = time.Now().Format("2006-01-02")
	}
	tglKStr := strings.TrimSpace(req.TglKembali)
	if tglKStr == "" {
		tglKStr = tglBStr
	}

	tglB, errB := time.Parse("2006-01-02", tglBStr)
	tglK, errK := time.Parse("2006-01-02", tglKStr)
	if errB != nil {
		tglB = time.Now()
	}
	if errK != nil {
		tglK = tglB
	}
	if tglK.Before(tglB) {
		tglK = tglB
	}

	// 2. Hitung durasi jam spesifik (DateTime Berangkat s/d Kembali)
	startDT, errStart := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", tglB.Format("2006-01-02"), jamB))
	endDT, errEnd := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", tglK.Format("2006-01-02"), jamK))
	durasiJam := 9.0
	if errStart == nil && errEnd == nil {
		diff := endDT.Sub(startDT).Hours()
		if diff > 0 {
			durasiJam = diff
		} else {
			durasiJam = 1.0
		}
	}

	// Durasi hari kalender
	durasiHari := int(tglK.Sub(tglB).Hours()/24) + 1
	if durasiHari < 1 {
		durasiHari = 1
	}

	// 3. Evaluasi Syarat Validasi (Jarak & Jam) dari Ketentuan Fasilitas (Pukul Rata Seluruh Rombongan)
	var hakMakanFactor float64 = 1.0
	var hakUangSakuAllowed bool = true
	var hakHotelAllowed bool = (durasiHari > 1)
	malamHotel := 0

	if durasiHari == 1 {
		if durasiJam < 5.0 {
			hakMakanFactor = 0.5 // Setengah hak uang makan jika kurang dari 5 jam
			hakUangSakuAllowed = false
		} else if durasiJam <= 8.0 {
			hakMakanFactor = 1.0
			hakUangSakuAllowed = false
		} else {
			hakMakanFactor = 1.0
			hakUangSakuAllowed = true
		}
		hakHotelAllowed = false
	} else {
		hakMakanFactor = 1.0
		hakUangSakuAllowed = true
		hakHotelAllowed = true
		malamHotel = durasiHari - 1
		if malamHotel < 1 {
			malamHotel = 1
		}
	}

	// 4. Kumpulkan Seluruh NomorID Personil (Pemohon Utama + Seluruh Partner)
	type PersonTarget struct {
		TargetID string
		RawText  string
	}
	var allTargets []PersonTarget
	seenIDs := make(map[string]bool)

	var pemohonNomorID string
	var pemohonNama string
	var pemohonKodeJabatan string

	if req.NomorID != "" {
		pemohonNomorID = strings.TrimSpace(req.NomorID)
	} else if idPemohon > 0 && s.userRepo != nil {
		if u, err := s.userRepo.FindByID(idPemohon); err == nil && u != nil {
			pemohonNomorID = strings.TrimSpace(u.NomorID)
			pemohonNama = strings.TrimSpace(u.NamaLengkap)
			pemohonKodeJabatan = strings.TrimSpace(u.KodeJabatan)
		}
	}

	if pemohonNomorID != "" {
		allTargets = append(allTargets, PersonTarget{TargetID: pemohonNomorID, RawText: pemohonNama})
		seenIDs[pemohonNomorID] = true
	} else {
		// Fallback pemohon jika nomor ID belum teridentifikasi
		allTargets = append(allTargets, PersonTarget{TargetID: "PEMOHON", RawText: pemohonNama})
	}

	for _, rekan := range req.RekanKerja {
		clean := strings.TrimSpace(rekan)
		if clean == "" {
			continue
		}
		// Ekstrak jika berformat "Nama (NomorID)"
		targetID := clean
		if lastOpen := strings.LastIndex(clean, "("); lastOpen != -1 {
			if lastClose := strings.LastIndex(clean, ")"); lastClose != -1 && lastClose > lastOpen {
				targetID = strings.TrimSpace(clean[lastOpen+1 : lastClose])
			}
		}
		if targetID != "" && !seenIDs[targetID] {
			seenIDs[targetID] = true
			allTargets = append(allTargets, PersonTarget{TargetID: targetID, RawText: clean})
		}
	}

	area := DetermineArea(req.KotaTujuan)

	// Ambil default tarif Area
	var defaultTarif models.PerjalananDinasTarif
	if err := config.DB.Where("Area = ?", area).First(&defaultTarif).Error; err != nil {
		_ = config.DB.Where("Area = 'Jawa'").First(&defaultTarif)
	}

	// 5. Query dan Susun Data Personil Rombongan Secara Lengkap
	var personilList []PersonilDinasTarif

	for idx, t := range allTargets {
		nama := t.RawText
		kodeJabatan := ""
		groupDinasLuar := ""
		nomorID := t.TargetID

		// Potong nama jika RawText berformat "Nama (ID) - Dept: ..."
		if strings.Contains(nama, "(") {
			nama = strings.TrimSpace(strings.Split(nama, "(")[0])
		}

		cleanTargetID := strings.TrimSpace(t.TargetID)
		cleanNama := strings.TrimSpace(nama)

		// Jika ini pemohon dan info pemohon sudah ada
		if idx == 0 && pemohonKodeJabatan != "" {
			kodeJabatan = pemohonKodeJabatan
		}

		var found bool

		// 1. Cari berdasarkan TargetID / NomorID terlebih dahulu (Prioritas Utama)
		if cleanTargetID != "" && !strings.EqualFold(cleanTargetID, "PEMOHON") {
			var k models.Karyawan
			if err := config.DB.Table("karyawan_copy1").Where("TRIM(REPLACE(REPLACE(NomorID, '\r', ''), '\n', '')) = ?", cleanTargetID).First(&k).Error; err == nil {
				nama = strings.TrimSpace(k.Nama)
				kodeJabatan = strings.TrimSpace(k.KodeJabatan)
				nomorID = strings.TrimSpace(k.NomorID)
				found = true
			} else {
				var u models.User
				if err := config.DB.Where("TRIM(REPLACE(REPLACE(nomor_id, '\r', ''), '\n', '')) = ?", cleanTargetID).First(&u).Error; err == nil {
					nama = strings.TrimSpace(u.NamaLengkap)
					kodeJabatan = strings.TrimSpace(u.KodeJabatan)
					nomorID = strings.TrimSpace(u.NomorID)
					found = true
				}
			}
		}

		// 2. Jika tidak ditemukan via ID, cari berdasarkan Nama (HANYA jika nama tidak kosong)
		if !found && cleanNama != "" && !strings.EqualFold(cleanNama, "PEMOHON") {
			var k models.Karyawan
			if err := config.DB.Table("karyawan_copy1").Where("TRIM(REPLACE(REPLACE(Nama, '\r', ''), '\n', '')) = ? OR TRIM(REPLACE(REPLACE(Nama, '\r', ''), '\n', '')) LIKE ?", cleanNama, "%"+cleanNama+"%").First(&k).Error; err == nil {
				nama = strings.TrimSpace(k.Nama)
				kodeJabatan = strings.TrimSpace(k.KodeJabatan)
				nomorID = strings.TrimSpace(k.NomorID)
				found = true
			} else {
				var u models.User
				if err := config.DB.Where("TRIM(REPLACE(REPLACE(nama_lengkap, '\r', ''), '\n', '')) = ? OR TRIM(REPLACE(REPLACE(nama_lengkap, '\r', ''), '\n', '')) LIKE ?", cleanNama, "%"+cleanNama+"%").First(&u).Error; err == nil {
					nama = strings.TrimSpace(u.NamaLengkap)
					kodeJabatan = strings.TrimSpace(u.KodeJabatan)
					nomorID = strings.TrimSpace(u.NomorID)
					found = true
				}
			}
		}

		// 3. Fallback jika idx == 0 (Pemohon) dan belum teridentifikasi dari query di atas
		if idx == 0 && (nama == "" || strings.EqualFold(nama, "PEMOHON")) {
			if idPemohon > 0 && s.userRepo != nil {
				if u, err := s.userRepo.FindByID(idPemohon); err == nil && u != nil {
					nama = strings.TrimSpace(u.NamaLengkap)
					nomorID = strings.TrimSpace(u.NomorID)
					if kodeJabatan == "" {
						kodeJabatan = strings.TrimSpace(u.KodeJabatan)
					}
				}
			}
		}

		if nama == "" || strings.EqualFold(nama, "PEMOHON") {
			if idx == 0 {
				nama = "Pemohon"
			} else {
				nama = t.TargetID
			}
		}

		// C. Cari GroupDinasLuar di tabel master jabatan
		if kodeJabatan != "" {
			var j models.Jabatan
			if err := config.DB.Where("TRIM(REPLACE(REPLACE(Kode, '\r', ''), '\n', '')) = ?", kodeJabatan).First(&j).Error; err == nil {
				groupDinasLuar = strings.TrimSpace(j.GroupDinasLuar)
			}
		}

		// Fallback pemetaan group dinas luar dari kode jabatan
		if groupDinasLuar == "" && kodeJabatan != "" {
			groupDinasLuar = MapKodeJabatanToGroup(kodeJabatan)
		}

		// D. Cari Tarif dari perjalanandinas_tarif berdasarkan GroupDinasLuar
		var tarif models.PerjalananDinasTarif
		tarifFound := false
		if groupDinasLuar != "" {
			if err := config.DB.Where("TRIM(Jabatan) = ? AND Area = ?", groupDinasLuar, area).First(&tarif).Error; err == nil {
				tarifFound = true
			}
		}
		if !tarifFound {
			tarif = defaultTarif
		}

		displayJabatan := groupDinasLuar
		if displayJabatan == "" {
			displayJabatan = "Operator - Karu"
		}

		personilList = append(personilList, PersonilDinasTarif{
			NomorID:        nomorID,
			Nama:           nama,
			KodeJabatan:    kodeJabatan,
			NamaJabatan:    displayJabatan,
			GroupDinasLuar: displayJabatan,
			Area:           area,
			MataUang:       "IDR",
			UangMakan:      tarif.UangMakan,
			UangSaku:       tarif.UangSaku,
			Hotel:          tarif.Hotel,
			Transport:      tarif.Transport,
			Tiket:          tarif.Tiket,
		})
		_ = idx
	}

	// 6. Loop Perhitungan Kalkulasi Per Kepala
	var plafonMakanTotal, plafonUangSakuTotal, plafonHotelTotal, plafonTransportTotal, plafonTiketTotal float64
	var personilBreakdown []string

	isHotel := false
	akom := strings.ToLower(strings.TrimSpace(req.Akomodasi))
	if strings.Contains(akom, "hotel") {
		isHotel = true
	} else if akom != "" && !strings.Contains(akom, "tidak") && !strings.Contains(akom, "mess") && !strings.Contains(akom, "pribadi") && !strings.Contains(akom, "tanpa") {
		isHotel = true
	}

	trans := strings.ToLower(strings.TrimSpace(req.Transportasi))
	isFreeTransport := strings.Contains(trans, "dinas") || strings.Contains(trans, "kantor") || strings.Contains(trans, "operasional") || strings.Contains(trans, "tanpa") || strings.Contains(trans, "tidak")

	var rincianPerKepala []dto.RincianPerKepalaDTO

	for idx, p := range personilList {
		rateMakan := p.UangMakan
		rateSaku := p.UangSaku
		rateHotel := p.Hotel
		rateTransport := p.Transport
		rateTiket := p.Tiket

		// Fallback tarif jika 0
		if rateMakan == 0 && defaultTarif.UangMakan > 0 {
			rateMakan = defaultTarif.UangMakan
		}
		if rateSaku == 0 && defaultTarif.UangSaku > 0 {
			rateSaku = defaultTarif.UangSaku
		}
		if rateHotel == 0 && defaultTarif.Hotel > 0 {
			rateHotel = defaultTarif.Hotel
		}
		if rateTransport == 0 && defaultTarif.Transport > 0 {
			rateTransport = defaultTarif.Transport
		}

		// A. Uang Makan per Kepala
		personMakan := rateMakan * hakMakanFactor * float64(durasiHari)
		plafonMakanTotal += personMakan

		// B. Uang Saku per Kepala
		personSaku := 0.0
		if hakUangSakuAllowed {
			personSaku = rateSaku * float64(durasiHari)
		}
		plafonUangSakuTotal += personSaku

		// C. Hotel per Kepala (Dihitung per kepala / kamar sendiri jika akomodasi hotel)
		personHotel := 0.0
		if hakHotelAllowed && isHotel {
			personHotel = rateHotel * float64(malamHotel)
		}
		plafonHotelTotal += personHotel

		// D. Transport per Kepala
		personTransport := 0.0
		if !isFreeTransport && trans != "" {
			personTransport = rateTransport
			if req.JarakKM > 0 {
				tarifPerKM := 1500.0
				if strings.Contains(trans, "motor") {
					tarifPerKM = 750.0
				}
				biayaKM := req.JarakKM * tarifPerKM
				if biayaKM > rateTransport {
					personTransport = biayaKM
				}
			}
		}
		plafonTransportTotal += personTransport

		// E. Tiket per Kepala
		personTiket := 0.0
		if !isFreeTransport && rateTiket > 0 {
			personTiket = rateTiket
		}
		plafonTiketTotal += personTiket

		role := "Pemohon"
		if idx > 0 {
			role = "Rombongan"
		}
		roleDisplay := p.GroupDinasLuar
		if roleDisplay == "" {
			roleDisplay = p.NamaJabatan
		}
		if roleDisplay == "" {
			roleDisplay = "Karyawan"
		}

		personSubtotal := personMakan + personSaku + personHotel

		rincianPerKepala = append(rincianPerKepala, dto.RincianPerKepalaDTO{
			NomorID:        p.NomorID,
			Nama:           p.Nama,
			Role:           role,
			KodeJabatan:    p.KodeJabatan,
			NamaJabatan:    p.NamaJabatan,
			GroupDinasLuar: roleDisplay,
			Makan:          personMakan,
			Saku:           personSaku,
			Hotel:          personHotel,
			Transport:      personTransport,
			Tiket:          personTiket,
			Subtotal:       personSubtotal,
		})

		personilBreakdown = append(personilBreakdown, fmt.Sprintf(
			"[%s] %s (%s): Makan %s, Saku %s, Hotel %s -> Subtotal: %s",
			role, p.Nama, roleDisplay,
			formatRupiahGo(personMakan), formatRupiahGo(personSaku),
			formatRupiahGo(personHotel),
			formatRupiahGo(personSubtotal),
		))
	}

	komponenBersama := dto.KomponenBersamaDTO{
		Transport: plafonTransportTotal,
		Tiket:     plafonTiketTotal,
		Total:     plafonTransportTotal + plafonTiketTotal,
	}

	// 7. Akumulasi Plafon Total dan Estimasi Biaya
	estimasiBiayaTotal := plafonMakanTotal + plafonUangSakuTotal + plafonHotelTotal + plafonTransportTotal + plafonTiketTotal
	totalUangHarian := plafonMakanTotal + plafonUangSakuTotal

	jumlahOrang := len(personilList)
	var jarakStr string
	if req.JarakKM > 0 {
		jarakStr = fmt.Sprintf(" | Jarak: %.1f KM", req.JarakKM)
	}

	rincianTeks := fmt.Sprintf(
		"Area: %s%s | %s-%s (%.1f Jam / %d Hari) | %d Orang | Plafon Total: %s (Makan: %s, Saku: %s, Hotel: %s, Trans: %s) | Rincian: %s",
		area, jarakStr, jamB, jamK, durasiJam, durasiHari, jumlahOrang,
		formatRupiahGo(estimasiBiayaTotal), formatRupiahGo(plafonMakanTotal), formatRupiahGo(plafonUangSakuTotal),
		formatRupiahGo(plafonHotelTotal), formatRupiahGo(plafonTransportTotal),
		strings.Join(personilBreakdown, "; "),
	)

	resp := &dto.KalkulasiBiayaResponse{
		TotalKeseluruhan:    estimasiBiayaTotal,
		RincianPerKepala:    rincianPerKepala,
		KomponenBersama:     komponenBersama,
		DurasiHari:          durasiHari,
		DurasiJam:           durasiJam,
		JarakKM:             req.JarakKM,
		JumlahOrang:         jumlahOrang,
		Area:                area,
		KategoriJabatan:     personilList[0].GroupDinasLuar,
		KodeJabatan:         personilList[0].KodeJabatan,
		MataUang:            "IDR",
		UangMakanPerHari:    plafonMakanTotal / float64(durasiHari*jumlahOrang),
		UangSakuPerHari:     plafonUangSakuTotal / float64(durasiHari*jumlahOrang),
		UangHarianPerHari:   (plafonMakanTotal + plafonUangSakuTotal) / float64(durasiHari*jumlahOrang),
		HotelPerMalam:       plafonHotelTotal / float64(jumlahOrang),
		TransportPerOrang:   plafonTransportTotal / float64(jumlahOrang),
		TarifTransportPerKM: 1500.0,

		PlafonMakan:     plafonMakanTotal,
		PlafonTransport: plafonTransportTotal,
		PlafonUangSaku:  plafonUangSakuTotal,
		PlafonTiket:     plafonTiketTotal,
		PlafonHotel:     plafonHotelTotal,
		PlafonTotal:     estimasiBiayaTotal,

		HakTransport: 1,
		HakUangMakan: 1,
		HakUangSaku:  1,
		HakTiket:     1,
		HakHotel:     1,

		TotalUangHarian:    totalUangHarian,
		TotalHotel:         plafonHotelTotal,
		TotalTransport:     plafonTransportTotal,
		EstimasiBiayaTotal: estimasiBiayaTotal,
		RincianTeks:        rincianTeks,
	}

	return resp, nil
}

func formatRupiahGo(n float64) string {
	in := strconv.FormatFloat(n, 'f', 0, 64)
	var out []rune
	l := len(in)
	for i, r := range in {
		out = append(out, r)
		if (l-i-1)%3 == 0 && i != l-1 {
			out = append(out, '.')
		}
	}
	return "Rp " + string(out)
}
