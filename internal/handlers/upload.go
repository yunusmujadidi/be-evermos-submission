package handlers

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UploadFile(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
    }
    
    ext := filepath.Ext(file.Filename)
    filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), "upload", ext)
    
    uploadPath := "./uploads/" + filename
    
    if err := c.SaveFile(file, uploadPath); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to save file"})
    }
    
    fileURL := fmt.Sprintf("/uploads/%s", filename)
    
    return c.JSON(fiber.Map{
        "message": "File uploaded successfully",
        "url":     fileURL,
        "filename": filename,
        "size":    file.Size,
    })
} 