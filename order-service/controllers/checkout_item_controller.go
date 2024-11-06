package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/order-service/dto"
	"github.com/ridhotamma/yourkasa/order-service/models"
	"gorm.io/gorm"
)

type CheckoutController struct {
	db *gorm.DB
}

func NewCheckoutController(db *gorm.DB) *CheckoutController {
	return &CheckoutController{db: db}
}

func (c *CheckoutController) AddToCart(ctx *gin.Context) {
	var input dto.CreateCheckoutItemDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID := ctx.GetString("customer_id")

	// Verify product exists and get its price
	var product models.Product
	if err := c.db.First(&product, input.ProductID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Get variant price if specified
	var price float64 = product.Price
	if input.VariantID != nil {
		var variant models.ProductVariant
		if err := c.db.First(&variant, input.VariantID).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
			return
		}
		price = variant.Price
	}

	// Validate and calculate addons
	var addonsData []byte
	if len(input.AddonsData) > 0 {
		addonsData, _ = json.Marshal(input.AddonsData)
	}

	checkoutItem := models.CheckoutItem{
		CustomerID: customerID,
		ProductID:  input.ProductID,
		VariantID:  input.VariantID,
		Quantity:   input.Quantity,
		Price:      price,
		AddonsData: string(addonsData),
		Notes:      input.Notes,
	}

	if err := c.db.Create(&checkoutItem).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Item added to cart", "id": checkoutItem.ID})
}

func (c *CheckoutController) UpdateCartItem(ctx *gin.Context) {
	id := ctx.Param("id")
	customerID := ctx.GetString("customer_id")

	var input dto.UpdateCheckoutItemDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var checkoutItem models.CheckoutItem
	if err := c.db.Where("id = ? AND customer_id = ?", id, customerID).First(&checkoutItem).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	updates := map[string]interface{}{}
	if input.Quantity > 0 {
		updates["quantity"] = input.Quantity
	}
	if len(input.AddonsData) > 0 {
		addonsData, _ := json.Marshal(input.AddonsData)
		updates["addons_data"] = string(addonsData)
	}
	if input.Notes != "" {
		updates["notes"] = input.Notes
	}
	if input.IsSelected != nil {
		updates["is_selected"] = *input.IsSelected
	}

	if err := c.db.Model(&checkoutItem).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cart item updated successfully"})
}

func (c *CheckoutController) RemoveFromCart(ctx *gin.Context) {
	id := ctx.Param("id")
	customerID := ctx.GetString("customer_id")

	result := c.db.Where("id = ? AND customer_id = ?", id, customerID).Delete(&models.CheckoutItem{})
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

func (c *CheckoutController) GetCart(ctx *gin.Context) {
	customerID := ctx.GetString("customer_id")

	var items []models.CheckoutItem
	if err := c.db.Where("customer_id = ?", customerID).
		Preload("Product").
		Preload("Variant").
		Find(&items).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}
