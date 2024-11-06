package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/product-service/dto"
	"github.com/ridhotamma/yourkasa/product-service/models"
	"gorm.io/gorm"
)

type AddonController struct {
	db *gorm.DB
}

func NewAddonController(db *gorm.DB) *AddonController {
	return &AddonController{db: db}
}

func (c *AddonController) Create(ctx *gin.Context) {
	var input dto.CreateAddonDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if SKU exists
	var existingAddon models.ProductAddon
	if err := c.db.Where("sku = ?", input.SKU).First(&existingAddon).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "SKU already exists"})
		return
	}

	addon := models.ProductAddon{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		SKU:         input.SKU,
		ImageURL:    input.ImageURL,
		Stock:       input.Stock,
		IsRequired:  input.IsRequired,
		MaxQuantity: input.MaxQuantity,
	}

	if err := c.db.Create(&addon).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create addon"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Addon created successfully", "id": addon.ID})
}

func (c *AddonController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var input dto.UpdateAddonDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var addon models.ProductAddon
	if err := c.db.First(&addon, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Addon not found"})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Price > 0 {
		updates["price"] = input.Price
	}
	if input.Stock >= 0 {
		updates["stock"] = input.Stock
	}
	if input.IsRequired != nil {
		updates["is_required"] = *input.IsRequired
	}
	if input.MaxQuantity != nil {
		updates["max_quantity"] = *input.MaxQuantity
	}

	if err := c.db.Model(&addon).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update addon"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Addon updated successfully"})
}

func (c *AddonController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.db.Delete(&models.ProductAddon{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete addon"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Addon deleted successfully"})
}

func (c *AddonController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var addon models.ProductAddon
	if err := c.db.Preload("Products").First(&addon, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Addon not found"})
		return
	}

	ctx.JSON(http.StatusOK, addon)
}

func (c *AddonController) GetByProductID(ctx *gin.Context) {
	productID := ctx.Param("productId")
	var addons []models.ProductAddon
	if err := c.db.Joins("JOIN product_addon_mappings ON product_addon_mappings.addon_id = product_addons.id").
		Where("product_addon_mappings.product_id = ?", productID).
		Find(&addons).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch addons"})
		return
	}

	ctx.JSON(http.StatusOK, addons)
}
