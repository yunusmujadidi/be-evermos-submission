package handlers

import (
	"be-evermos-submission/internal/database"
	"be-evermos-submission/internal/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
	}
	
	if err := database.DB.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Phone already exists"})
	}
	
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
		IsAdmin:  false,
	}
	
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}
	
	store := models.Store{
		UserID: user.ID,
		Name:   req.Name + "'s Shop",
	}
	database.DB.Create(&store)
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	
	return c.JSON(models.AuthResponse{
		Token: tokenString,
		User:  &user,
	})
}

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	
	return c.JSON(models.AuthResponse{
		Token: tokenString,
		User:  &user,
	})
}

func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	
	var user models.User
	if err := database.DB.Preload("Store").Preload("Addresses").First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	
	return c.JSON(user)
}

func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	
	database.DB.Save(&user)
	return c.JSON(user)
} 