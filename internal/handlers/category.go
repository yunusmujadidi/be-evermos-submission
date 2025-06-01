package handlers

import (
	"be-evermos-submission/internal/database"
	"be-evermos-submission/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetCategories(c *fiber.Ctx) error {
    var categories []models.Category
    
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    page, _ := strconv.Atoi(c.Query("page", "1"))
    offset := (page - 1) * limit
    
    database.DB.Limit(limit).Offset(offset).Find(&categories)
    
    return c.JSON(fiber.Map{
        "data": categories,
        "page": page,
        "limit": limit,
    })
}

func CreateCategory(c *fiber.Ctx) error {
    var category models.Category
    if err := c.BodyParser(&category); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    if err := database.DB.Create(&category).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create category"})
    }
    
    return c.Status(201).JSON(category)
}

func UpdateCategory(c *fiber.Ctx) error {
    id := c.Params("id")
    
    var category models.Category
    if err := database.DB.First(&category, id).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
    }
    
    if err := c.BodyParser(&category); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    database.DB.Save(&category)
    return c.JSON(category)
}

func DeleteCategory(c *fiber.Ctx) error {
    id := c.Params("id")
    
    if err := database.DB.Delete(&models.Category{}, id).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
    }
    
    return c.JSON(fiber.Map{"message": "Category deleted"})
} 