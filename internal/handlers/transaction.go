package handlers

import (
	"be-evermos-submission/internal/database"
	"be-evermos-submission/internal/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CreateTransactionRequest struct {
    AddressID     uint                         `json:"alamat_kirim"`
    PaymentMethod string                       `json:"metode_bayar"`
    Details       []CreateDetailTransactionRequest `json:"detail_trx"`
}

type CreateDetailTransactionRequest struct {
    ProductID uint `json:"id_produk"`
    Quantity  int  `json:"kuantitas"`
}

func GetMyTransactions(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    page, _ := strconv.Atoi(c.Query("page", "1"))
    offset := (page - 1) * limit
    
    var transactions []models.Transaction
    database.DB.Where("user_id = ?", userID).
        Preload("Details").
        Preload("Details.Product").
        Preload("Details.Store").
        Preload("Logs").
        Limit(limit).Offset(offset).
        Find(&transactions)
    
    return c.JSON(fiber.Map{
        "data": transactions,
        "page": page,
        "limit": limit,
    })
}

func CreateTransaction(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    
    var req CreateTransactionRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    var address models.Address
    if err := database.DB.Where("id = ? AND user_id = ?", req.AddressID, userID).First(&address).Error; err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid address"})
    }
    
    var totalPrice int
    var details []models.TransactionDetail
    
    for _, detail := range req.Details {
        var product models.Product
        if err := database.DB.Preload("Store").First(&product, detail.ProductID).Error; err != nil {
            return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Product %d not found", detail.ProductID)})
        }
        
        if product.Stock < detail.Quantity {
            return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Insufficient stock for product %s", product.Name)})
        }
        
        price, _ := strconv.Atoi(product.ConsumerPrice)
        detailTotal := price * detail.Quantity
        totalPrice += detailTotal
        
        details = append(details, models.TransactionDetail{
            ProductID:  detail.ProductID,
            StoreID:    product.StoreID,
            Quantity:   detail.Quantity,
            TotalPrice: detailTotal,
        })
    }
    
    invoiceCode := fmt.Sprintf("INV-%d-%d", userID, time.Now().Unix())
    
    transaction := models.Transaction{
        UserID:        userID,
        AddressID:     req.AddressID,
        TotalPrice:    totalPrice,
        Invoice:       invoiceCode,
        PaymentMethod: req.PaymentMethod,
    }
    
    if err := database.DB.Create(&transaction).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create transaction"})
    }
    
    for i, detail := range details {
        detail.TransactionID = transaction.ID
        if err := database.DB.Create(&detail).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to create transaction detail"})
        }
        
        var product models.Product
        database.DB.First(&product, detail.ProductID)
        stockBefore := product.Stock
        product.Stock -= detail.Quantity
        database.DB.Save(&product)
        
        productLog := models.ProductLog{
            TransactionID: transaction.ID,
            ProductName:   product.Name,
            ProductPrice:  detail.TotalPrice / detail.Quantity,
            Quantity:      detail.Quantity,
            StockBefore:   stockBefore,
            StockAfter:    product.Stock,
        }
        database.DB.Create(&productLog)
        
        details[i] = detail
    }
    
    database.DB.Preload("Details").Preload("Details.Product").Preload("Logs").First(&transaction, transaction.ID)
    
    return c.Status(201).JSON(transaction)
}

func GetTransactionByID(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    transactionID := c.Params("id")
    
    var transaction models.Transaction
    if err := database.DB.Where("id = ? AND user_id = ?", transactionID, userID).
        Preload("Details").
        Preload("Details.Product").
        Preload("Details.Store").
        Preload("Logs").
        First(&transaction).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
    }
    
    return c.JSON(transaction)
} 