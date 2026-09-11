package repository

import (
	"strings"

	"sipadin/models"

	"gorm.io/gorm"
)

// UserRepository mendefinisikan kontrak akses data untuk tabel users
type UserRepository interface {
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

// FindByIdentifier mencari user berdasarkan nama_lengkap, nomor_id, nik_pendek, atau role_sistem
func (r *userRepository) FindByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	clean := strings.TrimSpace(identifier)
	roleSearch := clean
	if strings.EqualFold(clean, "admin") {
		roleSearch = "HRGA"
	}

	likePattern := clean + "%"
	err := r.db.Where(
		"LOWER(nama_lengkap) = LOWER(?) OR LOWER(nama_lengkap) LIKE LOWER(?) OR nomor_id = ? OR nik_pendek = ? OR LOWER(role_sistem) = LOWER(?)",
		clean, likePattern, clean, clean, roleSearch,
	).
		Where("status_aktif = 1").
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Order("nama_lengkap ASC").Find(&users).Error
	return users, err
}

func (r *userRepository) FindByRole(role string) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("role_sistem = ? AND status_aktif = 1", role).Find(&users).Error
	return users, err
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id int) error {
	return r.db.Delete(&models.User{}, id).Error
}
