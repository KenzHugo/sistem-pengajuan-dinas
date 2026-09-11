package service

import (
	"errors"
	"log"
	"sipadin/config"
	"sipadin/dto"
	"sipadin/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService mendefinisikan kontrak business logic autentikasi
type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetActiveUsers() ([]dto.UserProfile, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// GetActiveUsers mengambil seluruh data user aktif langsung dari database
func (s *authService) GetActiveUsers() ([]dto.UserProfile, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var profiles []dto.UserProfile
	for _, u := range users {
		if u.StatusAktif == 1 {
			profiles = append(profiles, dto.UserProfile{
				IDUser:      u.IDUser,
				NomorID:     u.NomorID,
				NamaLengkap: u.NamaLengkap,
				RoleSistem:  u.RoleSistem,
				NoWhatsapp:  u.NoWhatsapp,
			})
		}
	}
	return profiles, nil
}

// Login memvalidasi kredensial dan mengembalikan JWT token
func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	log.Printf("[Auth] Upaya login dengan identifier: '%s'", req.Identifier)

	user, err := s.userRepo.FindByIdentifier(req.Identifier)
	if err != nil {
		log.Printf("[Auth] User tidak ditemukan untuk identifier: '%s'", req.Identifier)
		return nil, errors.New("nama lengkap atau password salah")
	}

	// Verifikasi password: coba bcrypt hash, fallback ke plain text jika data dimasukkan manual ke DB
	passwordValid := false
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err == nil {
		passwordValid = true
	} else if user.PasswordHash == req.Password {
		passwordValid = true
	}

	if !passwordValid {
		log.Printf("[Auth] Password tidak cocok untuk user: %s (%s)", user.NomorID, user.NamaLengkap)
		return nil, errors.New("nama lengkap atau password salah")
	}

	log.Printf("[Auth] Login berhasil untuk user: %s (%s, Role: %s)", user.NomorID, user.NamaLengkap, user.RoleSistem)

	// Buat JWT token
	claims := jwt.MapClaims{
		"id_user":      user.IDUser,
		"nama_lengkap": user.NamaLengkap,
		"role_sistem":  user.RoleSistem,
		"exp":          time.Now().Add(time.Duration(config.AppConfig.JWTExpireHours) * time.Hour).Unix(),
		"iat":          time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &dto.LoginResponse{
		Token: tokenString,
		User: dto.UserProfile{
			IDUser:      user.IDUser,
			NomorID:     user.NomorID,
			NamaLengkap: user.NamaLengkap,
			RoleSistem:  user.RoleSistem,
			NoWhatsapp:  user.NoWhatsapp,
		},
	}, nil
}
