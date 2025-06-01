package handlers

import (
	"be-evermos-submission/internal/database"
	"be-evermos-submission/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetMyAddresses(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    
    var addresses []models.Address
    database.DB.Where("user_id = ?", userID).Find(&addresses)
    
    return c.JSON(addresses)
}

func CreateAddress(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    
    var address models.Address
    if err := c.BodyParser(&address); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    address.UserID = userID
    
    if err := database.DB.Create(&address).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create address"})
    }
    
    return c.Status(201).JSON(address)
}

func GetAddressByID(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    addressID := c.Params("id")
    
    var address models.Address
    if err := database.DB.Where("id = ? AND user_id = ?", addressID, userID).First(&address).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Address not found"})
    }
    
    return c.JSON(address)
}

func UpdateAddress(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    addressID := c.Params("id")
    
    var address models.Address
    if err := database.DB.Where("id = ? AND user_id = ?", addressID, userID).First(&address).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Address not found"})
    }
    
    if err := c.BodyParser(&address); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    address.UserID = userID
    id, _ := strconv.Atoi(addressID)
    address.ID = uint(id)
    
    if err := database.DB.Save(&address).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to update address"})
    }
    
    return c.JSON(address)
}

func DeleteAddress(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    addressID := c.Params("id")
    
    if err := database.DB.Where("id = ? AND user_id = ?", addressID, userID).Delete(&models.Address{}).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Address not found"})
    }
    
    return c.JSON(fiber.Map{"message": "Address deleted"})
} 