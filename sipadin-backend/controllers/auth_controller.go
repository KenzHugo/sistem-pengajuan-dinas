package controllers

import (
	"sipadin-backend/config"
	"sipadin-backend/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengenkripsi kata sandi menggunakan bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash membandingkan kata sandi teks biasa dengan hash bcrypt
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	input.Email = strings.TrimSpace(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	if input.Email == "" || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Email dan kata sandi wajib diisi"})
	}

	var user models.User
	result := config.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Email atau kata sandi salah"})
	}

	// Cek kata sandi dengan bcrypt
	isMatch := CheckPasswordHash(input.Password, user.Password)

	// Fallback kompatibilitas: jika di DB masih teks biasa, verifikasi dan otomatis migrasikan ke Bcrypt
	if !isMatch && user.Password == input.Password {
		isMatch = true
		if newHashed, err := HashPassword(input.Password); err == nil {
			config.DB.Model(&user).Update("password", newHashed)
			user.Password = newHashed
		}
	}

	if !isMatch {
		return c.Status(401).JSON(fiber.Map{"message": "Email atau kata sandi salah"})
	}

	// Jangan kirim password hash ke response frontend untuk keamanan
	userResponse := user
	userResponse.Password = ""

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"user":    userResponse,
	})
}