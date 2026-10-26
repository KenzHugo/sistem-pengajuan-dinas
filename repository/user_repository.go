package repository

import (
	"strings"

	"sipadin/models"

	"gorm.io/gorm"
)

// UserRepository mendefinisikan kontrak akses data untuk autentikasi dan data user
type UserRepository interface {
	FindByUsernameOrName(identifier string) (*models.Karyawan, error)
	FindByIdentifier(identifier string) (*models.User, error)
	FindByID(id int) (*models.User, error)
	FindAll() ([]models.User, error)
	FindByRole(role string) ([]models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

const cleanKaryawanNamaExpr = "LTRIM(RTRIM(REPLACE(REPLACE(Nama, CHAR(13), ''), CHAR(10), '')))"
const cleanKaryawanNomorIDExpr = "LTRIM(RTRIM(REPLACE(REPLACE(NomorID, CHAR(13), ''), CHAR(10), '')))"
const cleanKaryawanLast5IDExpr = "RIGHT(" + cleanKaryawanNomorIDExpr + ", 5)"

// FindByUsernameOrName melakukan query langsung ke tabel adi_persons.dbo.karyawan berdasarkan 5 digit terakhir NomorID, NomorID lengkap, atau Nama
func (r *userRepository) FindByUsernameOrName(identifier string) (*models.Karyawan, error) {
	var karyawan models.Karyawan
	clean := strings.TrimSpace(identifier)

	err := r.db.Table("adi_persons.dbo.karyawan").
		Where("("+cleanKaryawanLast5IDExpr+" = ? OR "+cleanKaryawanNomorIDExpr+" = ? OR "+cleanKaryawanNomorIDExpr+" LIKE ? OR "+cleanKaryawanNamaExpr+" = ? OR LOWER("+cleanKaryawanNamaExpr+") = LOWER(?))",
			clean, clean, "%"+clean, clean, clean).
		Order("CASE WHEN Aktif = 1 THEN 0 ELSE 1 END, Nama ASC").
		First(&karyawan).Error

	if err != nil {
		return nil, err
	}

	karyawan.NomorID = cleanString(karyawan.NomorID)
	karyawan.Nama = cleanString(karyawan.Nama)
	karyawan.Password = cleanString(karyawan.Password)
	return &karyawan, nil
}

// FindByIdentifier mencocokkan kredensial langsung ke tabel adi_persons.dbo.karyawan dan memetakan ke struct User sistem
func (r *userRepository) FindByIdentifier(identifier string) (*models.User, error) {
	karyawan, err := r.FindByUsernameOrName(identifier)
	if err != nil {
		return nil, err
	}

	// Menentukan role sistem berdasarkan KodeJabatan dan KodeDepartemen
	role := "Karyawan"
	clean := strings.TrimSpace(identifier)
	jabatanUpper := strings.ToUpper(strings.TrimSpace(karyawan.KodeJabatan))
	deptUpper := strings.ToUpper(strings.TrimSpace(karyawan.KodeDepartemen))

	if strings.EqualFold(clean, "admin") {
		role = "HRGA"
	} else if jabatanUpper == "31" || jabatanUpper == "30" || strings.Contains(jabatanUpper, "DIR") {
		role = "Direksi"
	} else if jabatanUpper == "29" || jabatanUpper == "28" || strings.Contains(jabatanUpper, "MGR") || strings.Contains(jabatanUpper, "MGT") || deptUpper == "MGT" {
		role = "Manager"
	} else if deptUpper == "HRGA" || deptUpper == "HRD" || deptUpper == "GA" {
		role = "HRGA"
	} else if deptUpper == "FIN" || deptUpper == "FA" || deptUpper == "ACC" {
		role = "Finance"
	}

	// Hitung nik_pendek (5 digit terakhir NomorID)
	cleanNomorID := cleanString(karyawan.NomorID)
	nikPendek := cleanNomorID
	if len(cleanNomorID) >= 5 {
		nikPendek = cleanNomorID[len(cleanNomorID)-5:]
	}

	// Sinkronisasi/Cek record di tabel users untuk integritas relasi ID SPPD
	var existingUser models.User
	errUser := r.db.Where("LTRIM(RTRIM(REPLACE(REPLACE(nomor_id, CHAR(13), ''), CHAR(10), ''))) = ?", karyawan.NomorID).First(&existingUser).Error

	if errUser == nil {
		existingUser.NamaLengkap = karyawan.Nama
		existingUser.NikPendek = nikPendek
		existingUser.KodeDepartemen = karyawan.KodeDepartemen
		existingUser.KodeJabatan = karyawan.KodeJabatan
		existingUser.RoleSistem = role
		if karyawan.Password != "" {
			existingUser.PasswordHash = karyawan.Password
		}
		existingUser.StatusAktif = 1
		_ = r.db.Save(&existingUser).Error
		return &existingUser, nil
	}

	// Buat record User baru jika belum ada di tabel users
	newUser := models.User{
		NomorID:        karyawan.NomorID,
		NikPendek:      nikPendek,
		NamaLengkap:    karyawan.Nama,
		KodeDepartemen: karyawan.KodeDepartemen,
		KodeJabatan:    karyawan.KodeJabatan,
		RoleSistem:     role,
		NoWhatsapp:     "",
		PasswordHash:   karyawan.Password,
		StatusAktif:    1,
	}

	_ = r.db.Create(&newUser).Error
	return &newUser, nil
}

func (r *userRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	user.NomorID = cleanString(user.NomorID)
	user.NamaLengkap = cleanString(user.NamaLengkap)
	return &user, nil
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var listKaryawan []models.Karyawan
	// err := r.db.Table("karyawan_copy1").Order("Nama ASC").Find(&listKaryawan).Error
	err := r.db.Table("adi_persons.dbo.karyawan").Order("Nama ASC").Find(&listKaryawan).Error
	if err != nil {
		var users []models.User
		errUsers := r.db.Order("nama_lengkap ASC").Find(&users).Error
		return users, errUsers
	}

	var users []models.User
	for i, k := range listKaryawan {
		role := "Karyawan"
		j := strings.ToUpper(strings.TrimSpace(k.KodeJabatan))
		d := strings.ToUpper(strings.TrimSpace(k.KodeDepartemen))
		if j == "31" || j == "30" || strings.Contains(j, "DIR") {
			role = "Direksi"
		} else if j == "29" || j == "28" || strings.Contains(j, "MGR") || strings.Contains(j, "MGT") || d == "MGT" {
			role = "Manager"
		} else if d == "HRGA" || d == "HRD" || d == "GA" {
			role = "HRGA"
		} else if d == "FIN" || d == "FA" || d == "ACC" {
			role = "Finance"
		}

		cID := cleanString(k.NomorID)
		nikPendek := cID
		if len(cID) >= 5 {
			nikPendek = cID[len(cID)-5:]
		}

		users = append(users, models.User{
			IDUser:         i + 1,
			NomorID:        cID,
			NikPendek:      nikPendek,
			NamaLengkap:    cleanString(k.Nama),
			KodeDepartemen: k.KodeDepartemen,
			KodeJabatan:    k.KodeJabatan,
			RoleSistem:     role,
			StatusAktif:    1,
		})
	}

	return users, nil
}

func (r *userRepository) FindByRole(role string) ([]models.User, error) {
	allUsers, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	var filtered []models.User
	for _, u := range allUsers {
		if strings.EqualFold(u.RoleSistem, role) {
			filtered = append(filtered, u)
		}
	}
	return filtered, nil
}

func (r *userRepository) Create(user *models.User) error {
	user.NomorID = cleanString(user.NomorID)
	user.NamaLengkap = cleanString(user.NamaLengkap)
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	user.NomorID = cleanString(user.NomorID)
	user.NamaLengkap = cleanString(user.NamaLengkap)
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id int) error {
	return r.db.Delete(&models.User{}, id).Error
}

func cleanString(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return strings.TrimSpace(s)
}
