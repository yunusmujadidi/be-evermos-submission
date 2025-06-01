package handlers

import (
	"be-evermos-submission/internal/database"
	"be-evermos-submission/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetProducts(c *fiber.Ctx) error {
    var products []models.Product
    
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    page, _ := strconv.Atoi(c.Query("page", "1"))
    offset := (page - 1) * limit
    
    query := database.DB.Preload("Store").Preload("Category").Preload("Photos")
    
    if category := c.Query("category"); category != "" {
        query = query.Where("category_id = ?", category)
    }
    
    if search := c.Query("search"); search != "" {
        query = query.Where("name LIKE ?", "%"+search+"%")
    }
    
    query.Limit(limit).Offset(offset).Find(&products)
    
    return c.JSON(fiber.Map{
        "data": products,
        "page": page,
        "limit": limit,
    })
}

func GetMyProducts(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    
    var store models.Store
    if err := database.DB.Where("user_id = ?", userID).First(&store).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Shop not found"})
    }
    
    var products []models.Product
    database.DB.Where("store_id = ?", store.ID).Preload("Category").Preload("Photos").Find(&products)
    
    return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    
    var store models.Store
    if err := database.DB.Where("user_id = ?", userID).First(&store).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Shop not found"})
    }
    
    var product models.Product
    if err := c.BodyParser(&product); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    product.StoreID = store.ID
    
    if err := database.DB.Create(&product).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create product"})
    }
    
    return c.Status(201).JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    productID := c.Params("id")
    
    var store models.Store
    if err := database.DB.Where("user_id = ?", userID).First(&store).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Shop not found"})
    }
    
    var product models.Product
    if err := database.DB.Where("id = ? AND store_id = ?", productID, store.ID).First(&product).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
    }
    
    if err := c.BodyParser(&product); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    database.DB.Save(&product)
    return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    productID := c.Params("id")
    
    var store models.Store
    if err := database.DB.Where("user_id = ?", userID).First(&store).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Shop not found"})
    }
    
    if err := database.DB.Where("id = ? AND store_id = ?", productID, store.ID).Delete(&models.Product{}).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
    }
    
    return c.JSON(fiber.Map{"message": "Product deleted"})
} 