package repository

import (
	"strings"

	"sipadin/dto"
	"sipadin/models"

	"gorm.io/gorm"
)

// KaryawanRepository mendefinisikan kontrak akses data untuk tabel karyawan_copy1
type KaryawanRepository interface {
	FindAll(filter dto.KaryawanQuery) ([]models.Karyawan, int64, error)
	FindByID(nomorID string) (*models.Karyawan, error)
	Create(karyawan *models.Karyawan) error
	Update(nomorID string, karyawan *models.Karyawan) error
	Delete(nomorID string) error
	GetDepartemenList() ([]string, error)
	GetStats() (*dto.KaryawanStatsResponse, error)
}

type karyawanRepository struct {
	db *gorm.DB
}

// NewKaryawanRepository membuat instance baru KaryawanRepository
func NewKaryawanRepository(db *gorm.DB) KaryawanRepository {
	return &karyawanRepository{db: db}
}

// cleanIDQuery ekspresi SQL untuk membersihkan spasi dan karakter newline \r \n pada kolom NomorID
const cleanIDQuery = "TRIM(REPLACE(REPLACE(NomorID, '\r', ''), '\n', ''))"

// FindAll mengambil daftar data karyawan dengan filter pencarian, departemen, status aktif, dan pagination
func (r *karyawanRepository) FindAll(filter dto.KaryawanQuery) ([]models.Karyawan, int64, error) {
	var list []models.Karyawan
	var total int64

	query := r.db.Model(&models.Karyawan{})

	// Filter Search keyword (Nama, NomorID terbersihkan, Job, Kota, KodeDepartemen)
	if search := strings.TrimSpace(filter.Search); search != "" {
		likePattern := "%" + search + "%"
		query = query.Where(
			"Nama LIKE ? OR "+cleanIDQuery+" LIKE ? OR Job LIKE ? OR Kota LIKE ? OR KodeDepartemen LIKE ? OR KodeJabatan LIKE ?",
			likePattern, likePattern, likePattern, likePattern, likePattern, likePattern,
		)
	}

	// Filter Kode Departemen
	if dept := strings.TrimSpace(filter.KodeDepartemen); dept != "" {
		query = query.Where("KodeDepartemen = ?", dept)
	}

	// Filter Status Aktif ("1" = Aktif, "0" = Non-aktif)
	if filter.Aktif != "" {
		if filter.Aktif == "1" {
			query = query.Where("Aktif = 1")
		} else if filter.Aktif == "0" {
			query = query.Where("Aktif = 0")
		}
	}

	// Hitung total data sesuai filter
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	orderClause := "Nama ASC"
	if filter.SortBy != "" {
		orderDir := "ASC"
		if strings.ToUpper(filter.Order) == "DESC" {
			orderDir = "DESC"
		}
		switch filter.SortBy {
		case "nomor_id", "NomorID":
			orderClause = "NomorID " + orderDir
		case "nama", "Nama":
			orderClause = "Nama " + orderDir
		case "kode_departemen", "KodeDepartemen":
			orderClause = "KodeDepartemen " + orderDir
		case "job", "Job":
			orderClause = "Job " + orderDir
		case "kota", "Kota":
			orderClause = "Kota " + orderDir
		case "tanggal_lahir", "TanggalLahir":
			orderClause = "TanggalLahir " + orderDir
		case "aktif", "Aktif":
			orderClause = "Aktif " + orderDir
		}
	}
	query = query.Order(orderClause)

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 15 // Default limit per page
	}
	if limit > 200 {
		limit = 200
	}
	offset := (page - 1) * limit

	err := query.Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	// Clean NomorID di memori untuk data yang mungkin memiliki trailing \r\n
	for i := range list {
		list[i].NomorID = cleanStringID(list[i].NomorID)
	}

	return list, total, nil
}

// FindByID mencari karyawan berdasarkan NomorID dengan proteksi data cleaning query SQL
func (r *karyawanRepository) FindByID(nomorID string) (*models.Karyawan, error) {
	var karyawan models.Karyawan
	cleaned := cleanStringID(nomorID)

	err := r.db.Where(cleanIDQuery+" = ?", cleaned).First(&karyawan).Error
	if err != nil {
		return nil, err
	}

	karyawan.NomorID = cleanStringID(karyawan.NomorID)
	return &karyawan, nil
}

// Create menyimpan data karyawan baru ke database
func (r *karyawanRepository) Create(karyawan *models.Karyawan) error {
	karyawan.NomorID = cleanStringID(karyawan.NomorID)
	return r.db.Create(karyawan).Error
}

// Update memperbarui data karyawan berdasarkan NomorID dengan data cleaning
func (r *karyawanRepository) Update(nomorID string, karyawan *models.Karyawan) error {
	cleaned := cleanStringID(nomorID)
	return r.db.Model(&models.Karyawan{}).
		Where(cleanIDQuery+" = ?", cleaned).
		Updates(map[string]interface{}{
			"Nama":           karyawan.Nama,
			"KodeDepartemen": karyawan.KodeDepartemen,
			"KodeJabatan":    karyawan.KodeJabatan,
			"Job":            karyawan.Job,
			"Kota":           karyawan.Kota,
			"TanggalLahir":   karyawan.TanggalLahir,
			"Aktif":          karyawan.Aktif,
		}).Error
}

// Delete menghapus data karyawan berdasarkan NomorID dengan data cleaning
func (r *karyawanRepository) Delete(nomorID string) error {
	cleaned := cleanStringID(nomorID)
	return r.db.Where(cleanIDQuery+" = ?", cleaned).Delete(&models.Karyawan{}).Error
}

// GetDepartemenList mengambil daftar unik kode departemen yang terdaftar
func (r *karyawanRepository) GetDepartemenList() ([]string, error) {
	var depts []string
	err := r.db.Model(&models.Karyawan{}).
		Where("KodeDepartemen IS NOT NULL AND KodeDepartemen != ''").
		Distinct("KodeDepartemen").
		Order("KodeDepartemen ASC").
		Pluck("KodeDepartemen", &depts).Error
	return depts, err
}

// GetStats menghitung metrik ringkasan karyawan (total, aktif, non-aktif, departemen)
func (r *karyawanRepository) GetStats() (*dto.KaryawanStatsResponse, error) {
	var total, aktif, nonAktif, totalDept int64

	if err := r.db.Model(&models.Karyawan{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.Karyawan{}).Where("Aktif = 1").Count(&aktif).Error; err != nil {
		return nil, err
	}

	nonAktif = total - aktif

	if err := r.db.Model(&models.Karyawan{}).
		Where("KodeDepartemen IS NOT NULL AND KodeDepartemen != ''").
		Distinct("KodeDepartemen").
		Count(&totalDept).Error; err != nil {
		return nil, err
	}

	return &dto.KaryawanStatsResponse{
		TotalKaryawan:   total,
		KaryawanAktif:   aktif,
		KaryawanNonAktif: nonAktif,
		TotalDepartemen: totalDept,
	}, nil
}

// cleanStringID membersihkan karakter whitespace, \r, dan \n dari string ID
func cleanStringID(id string) string {
	id = strings.ReplaceAll(id, "\r", "")
	id = strings.ReplaceAll(id, "\n", "")
	return strings.TrimSpace(id)
}
