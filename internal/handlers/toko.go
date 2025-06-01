package handlers

import (
	"be-evermos-submission/internal/database"
	"be-evermos-submission/internal/models"

	"github.com/gofiber/fiber/v2"
)

func GetMyStore(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    
    var store models.Store
    if err := database.DB.Where("user_id = ?", userID).First(&store).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Store not found"})
    }
    
    return c.JSON(store)
}

func UpdateMyStore(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    
    var store models.Store
    if err := database.DB.Where("user_id = ?", userID).First(&store).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Store not found"})
    }
    
    if err := c.BodyParser(&store); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    store.UserID = userID
    
    if err := database.DB.Save(&store).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to update store"})
    }
    
    return c.JSON(store)
} 