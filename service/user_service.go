package service

import (
	"errors"
	"sipadin/dto"
	"sipadin/models"
	"sipadin/repository"

	"golang.org/x/crypto/bcrypt"
)

// UserService mendefinisikan business logic manajemen user
type UserService interface {
	GetAll() ([]models.User, error)
	GetByID(id int) (*models.User, error)
	Create(req dto.CreateUserRequest) (*models.User, error)
	Update(id int, req dto.UpdateUserRequest) (*models.User, error)
	Delete(id int) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *userService) GetByID(id int) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return user, nil
}

func (s *userService) Create(req dto.CreateUserRequest) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}

	statusAktif := req.StatusAktif
	if statusAktif == 0 {
		statusAktif = 1 // default aktif
	}

	user := &models.User{
		NomorID:        req.NomorID,
		NikPendek:      req.NikPendek,
		NamaLengkap:    req.NamaLengkap,
		KodeDepartemen: req.KodeDepartemen,
		KodeJabatan:    req.KodeJabatan,
		RoleSistem:     req.RoleSistem,
		NoWhatsapp:     req.NoWhatsapp,
		PasswordHash:   string(hash),
		StatusAktif:    statusAktif,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("gagal menyimpan user")
	}
	return user, nil
}

func (s *userService) Update(id int, req dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if req.NikPendek != "" {
		user.NikPendek = req.NikPendek
	}
	if req.NamaLengkap != "" {
		user.NamaLengkap = req.NamaLengkap
	}
	if req.KodeDepartemen != "" {
		user.KodeDepartemen = req.KodeDepartemen
	}
	if req.KodeJabatan != "" {
		user.KodeJabatan = req.KodeJabatan
	}
	if req.RoleSistem != "" {
		user.RoleSistem = req.RoleSistem
	}
	if req.NoWhatsapp != "" {
		user.NoWhatsapp = req.NoWhatsapp
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("gagal mengenkripsi password")
		}
		user.PasswordHash = string(hash)
	}
	if req.StatusAktif != nil {
		user.StatusAktif = *req.StatusAktif
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("gagal mengupdate user")
	}
	return user, nil
}

func (s *userService) Delete(id int) error {
	if _, err := s.userRepo.FindByID(id); err != nil {
		return errors.New("user tidak ditemukan")
	}
	return s.userRepo.Delete(id)
}
