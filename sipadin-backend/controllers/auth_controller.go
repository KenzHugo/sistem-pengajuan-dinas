package controllers

import (
	"sipadin-backend/config"
	"sipadin-backend/models"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	var user models.User
	result := config.DB.Where("email = ? AND password = ?", input.Email, input.Password).First(&user)
	if result.Error != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Email atau password salah"})
	}

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"user":    user,
	})
}