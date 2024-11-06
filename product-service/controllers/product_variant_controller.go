package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/product-service/dto"
	"github.com/ridhotamma/yourkasa/product-service/models"
	"gorm.io/gorm"
)

type VariantController struct {
	db *gorm.DB
}

func NewVariantController(db *gorm.DB) *VariantController {
	return &VariantController{db: db}
}

func (c *VariantController) Create(ctx *gin.Context) {
	var input dto.CreateVariantDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if product exists
	var product models.Product
	if err := c.db.First(&product, input.ProductID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	// Check if SKU exists
	var existingVariant models.ProductVariant
	if err := c.db.Where("sku = ?", input.SKU).First(&existingVariant).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "SKU already exists"})
		return
	}

	variant := models.ProductVariant{
		ProductID:  input.ProductID,
		Name:       input.Name,
		SKU:        input.SKU,
		Price:      input.Price,
		Stock:      input.Stock,
		ImageURL:   input.ImageURL,
		Attributes: input.Attributes,
		IsDefault:  input.IsDefault,
		Weight:     input.Weight,
		Dimensions: input.Dimensions,
		BarCode:    input.BarCode,
	}

	if err := c.db.Create(&variant).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create variant"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Variant created successfully", "id": variant.ID})
}

func (c *VariantController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var input dto.UpdateVariantDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var variant models.ProductVariant
	if err := c.db.First(&variant, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Price > 0 {
		updates["price"] = input.Price
	}
	if input.Stock >= 0 {
		updates["stock"] = input.Stock
	}
	if input.ImageURL != "" {
		updates["image_url"] = input.ImageURL
	}
	if input.Attributes != "" {
		updates["attributes"] = input.Attributes
	}
	if input.IsDefault != nil {
		updates["is_default"] = *input.IsDefault
	}

	if err := c.db.Model(&variant).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update variant"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Variant updated successfully"})
}

func (c *VariantController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.db.Delete(&models.ProductVariant{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete variant"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Variant deleted successfully"})
}

func (c *VariantController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var variant models.ProductVariant
	if err := c.db.Preload("Product").First(&variant, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	ctx.JSON(http.StatusOK, variant)
}

func (c *VariantController) GetByProductID(ctx *gin.Context) {
	productID := ctx.Param("productId")
	var variants []models.ProductVariant
	if err := c.db.Where("product_id = ?", productID).Find(&variants).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch variants"})
		return
	}

	ctx.JSON(http.StatusOK, variants)
}
