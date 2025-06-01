package main

import (
	"be-evermos-submission/internal/database"
	"be-evermos-submission/internal/handlers"
	"be-evermos-submission/internal/middleware"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")
	os.MkdirAll("uploads", 0755)
	
	database.Connect()
	
	app := fiber.New()
	
	app.Use(cors.New())
	app.Use(logger.New())
	
	app.Static("/uploads", "./uploads")
	
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	
	auth := app.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
	
	api := app.Group("/api", middleware.Auth)
	api.Get("/profile", handlers.GetProfile)
	api.Put("/profile", handlers.UpdateProfile)
	
	api.Get("/categories", handlers.GetCategories)
	admin := api.Group("/admin", middleware.AdminOnly)
	admin.Post("/categories", handlers.CreateCategory)
	admin.Put("/categories/:id", handlers.UpdateCategory)
	admin.Delete("/categories/:id", handlers.DeleteCategory)
	
	api.Get("/products", handlers.GetProducts)
	api.Get("/my-products", handlers.GetMyProducts)
	api.Post("/products", handlers.CreateProduct)
	api.Put("/products/:id", handlers.UpdateProduct)
	api.Delete("/products/:id", handlers.DeleteProduct)
	
	api.Get("/my-store", handlers.GetMyStore)
	api.Put("/my-store", handlers.UpdateMyStore)
	
	api.Get("/my-addresses", handlers.GetMyAddresses)
	api.Post("/my-addresses", handlers.CreateAddress)
	api.Get("/my-addresses/:id", handlers.GetAddressByID)
	api.Put("/my-addresses/:id", handlers.UpdateAddress)
	api.Delete("/my-addresses/:id", handlers.DeleteAddress)
	
	api.Get("/my-transactions", handlers.GetMyTransactions)
	api.Post("/my-transactions", handlers.CreateTransaction)
	api.Get("/my-transactions/:id", handlers.GetTransactionByID)
	
	api.Post("/upload", handlers.UploadFile)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
} 