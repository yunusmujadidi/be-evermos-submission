package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Auth(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"error": "No token provided"})
	}
	
	token = strings.Replace(token, "Bearer ", "", 1)
	
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	
	if err != nil || !parsed.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}
	
	claims := parsed.Claims.(jwt.MapClaims)
	c.Locals("user_id", uint(claims["user_id"].(float64)))
	c.Locals("is_admin", claims["is_admin"].(bool))
	
	return c.Next()
}

func AdminOnly(c *fiber.Ctx) error {
	isAdmin := c.Locals("is_admin").(bool)
	if !isAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only"})
	}
	return c.Next()
} 