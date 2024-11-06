package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/order-service/dto"
	"github.com/ridhotamma/yourkasa/order-service/models"
	"gorm.io/gorm"
)

type OrderController struct {
	db *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{db: db}
}

func (c *OrderController) Create(ctx *gin.Context) {
	var input dto.CreateOrderDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID := ctx.GetString("customer_id")

	// Get selected cart items
	var cartItems []models.CheckoutItem
	if err := c.db.Where("customer_id = ? AND is_selected = ?", customerID, true).
		Preload("Product").
		Preload("Variant").
		Find(&cartItems).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	if len(cartItems) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No items selected for checkout"})
		return
	}

	// Start transaction
	tx := c.db.Begin()

	// Create order
	orderNumber := fmt.Sprintf("ORD-%d-%s", time.Now().Unix(), customerID[:8])
	var subtotal float64
	var orderItems []models.OrderItem

	for _, item := range cartItems {
		// Calculate item total
		itemTotal := item.Price * float64(item.Quantity)

		// Add addons cost if any
		if item.AddonsData != "" {
			var addons []dto.CheckoutAddonDTO
			json.Unmarshal([]byte(item.AddonsData), &addons)
			// Calculate addons cost (you'll need to fetch addon prices)
		}

		orderItem := models.OrderItem{
			ProductID:   item.ProductID,
			VariantID:   item.VariantID,
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			Price:       item.Price,
			SubTotal:    itemTotal,
			AddonsData:  item.AddonsData,
		}
		if item.Variant != nil {
			orderItem.VariantName = item.Variant.Name
		}

		orderItems = append(orderItems, orderItem)
		subtotal += itemTotal
	}

	// Calculate tax and shipping (implement your business logic)
	taxAmount := subtotal * 0.1   // 10% tax example
	shippingAmount := float64(10) // Fixed shipping example

	order := models.Order{
		OrderNumber:     orderNumber,
		CustomerID:      customerID,
		Status:          "pending",
		SubtotalAmount:  subtotal,
		TaxAmount:       taxAmount,
		ShippingAmount:  shippingAmount,
		TotalAmount:     subtotal + taxAmount + shippingAmount,
		ShippingAddress: input.ShippingAddress,
		BillingAddress:  input.BillingAddress,
		PaymentMethod:   input.PaymentMethod,
		Notes:           input.Notes,
		OrderItems:      orderItems,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Clear cart items
	if err := tx.Where("id IN ?", getCartItemIDs(cartItems)).Delete(&models.CheckoutItem{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	tx.Commit()
	ctx.JSON(http.StatusCreated, gin.H{"message": "Order created successfully", "orderNumber": order.OrderNumber})
}

func (c *OrderController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	customerID := ctx.GetString("customer_id")

	var order models.Order
	if err := c.db.Where("id = ? AND customer_id = ?", id, customerID).
		Preload("OrderItems").
		First(&order).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (c *OrderController) List(ctx *gin.Context) {
	customerID := ctx.GetString("customer_id")

	var orders []models.Order
	if err := c.db.Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (c *OrderController) Cancel(ctx *gin.Context) {
	id := ctx.Param("id")
	customerID := ctx.GetString("customer_id")

	var order models.Order
	if err := c.db.Where("id = ? AND customer_id = ?", id, customerID).First(&order).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status != "pending" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only pending orders can be cancelled"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":      "cancelled",
		"canceled_at": now,
	}

	if err := c.db.Model(&order).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

// Helper functions
func getCartItemIDs(items []models.CheckoutItem) []uint {
	ids := make([]uint, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return ids
}
