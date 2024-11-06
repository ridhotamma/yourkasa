package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/product-service/dto"
	"github.com/ridhotamma/yourkasa/product-service/models"
	"gorm.io/gorm"
)

type ProductController struct {
	db *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{db: db}
}

func (c *ProductController) Create(ctx *gin.Context) {
	var input dto.CreateProductDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if SKU already exists
	var existingProduct models.Product
	if err := c.db.Where("sku = ?", input.SKU).First(&existingProduct).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "SKU already exists"})
		return
	}

	// Verify category exists
	if !categoryExists(c.db, input.CategoryID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}

	// Verify group exists if provided
	if input.GroupID != nil && !groupExists(c.db, *input.GroupID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Group not found"})
		return
	}

	product := models.Product{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		Price:            input.Price,
		Stock:            input.Stock,
		CategoryID:       input.CategoryID,
		GroupID:          input.GroupID,
		SKU:              input.SKU,
		ImageURL:         input.ImageURL,
		Weight:           input.Weight,
		Dimensions:       input.Dimensions,
		Tags:             input.Tags,
		Status:           input.Status,
	}

	// Start transaction
	tx := c.db.Begin()
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Associate addons if provided
	if len(input.AddonIDs) > 0 {
		var addons []models.ProductAddon
		if err := tx.Find(&addons, input.AddonIDs).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "One or more addons not found"})
			return
		}

		if err := tx.Model(&product).Association("Addons").Append(&addons); err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate addons"})
			return
		}
	}

	tx.Commit()
	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "id": product.ID})
}

func (c *ProductController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var input dto.UpdateProductDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := c.db.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.ShortDescription != "" {
		updates["short_description"] = input.ShortDescription
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
	if input.CategoryID != nil {
		if !categoryExists(c.db, *input.CategoryID) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
		updates["category_id"] = *input.CategoryID
	}
	if input.ImageURL != "" {
		updates["image_url"] = input.ImageURL
	}
	if input.Weight > 0 {
		updates["weight"] = input.Weight
	}
	if input.Dimensions != "" {
		updates["dimensions"] = input.Dimensions
	}
	if input.Tags != "" {
		updates["tags"] = input.Tags
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if input.GroupID != nil {
		if !groupExists(c.db, *input.GroupID) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Group not found"})
			return
		}
		updates["group_id"] = *input.GroupID
	}

	tx := c.db.Begin()

	if err := tx.Model(&product).Updates(updates).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Update addons if provided
	if len(input.AddonIDs) > 0 {
		if err := tx.Model(&product).Association("Addons").Clear(); err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear existing addons"})
			return
		}

		var addons []models.ProductAddon
		if err := tx.Find(&addons, input.AddonIDs).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "One or more addons not found"})
			return
		}

		if err := tx.Model(&product).Association("Addons").Append(&addons); err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate addons"})
			return
		}
	}

	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func (c *ProductController) List(ctx *gin.Context) {
	var products []models.Product
	query := c.db.Model(&models.Product{})

	// Apply filters
	if category := ctx.Query("category"); category != "" {
		query = query.Where("category_id = ?", category)
	}

	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if group := ctx.Query("group"); group != "" {
		query = query.Where("group_id = ?", group)
	}

	// Preload relationships
	query = query.Preload("Category").
		Preload("Group").
		Preload("Variants").
		Preload("Addons")

	if err := query.Find(&products).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	var productList []dto.ProductListDTO
	for _, product := range products {
		productList = append(productList, dto.ProductListDTO{
			ID:               product.ID,
			Name:             product.Name,
			ShortDescription: product.ShortDescription,
			Price:            product.Price,
			Stock:            product.Stock,
			CategoryID:       product.CategoryID,
			CategoryName:     product.Category.Name,
			SKU:              product.SKU,
			Status:           product.Status,
			VariantCount:     len(product.Variants),
		})
	}

	ctx.JSON(http.StatusOK, productList)
}

func (c *ProductController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.Product

	if err := c.db.Preload("Category").
		Preload("Group").
		Preload("Variants").
		Preload("Addons").
		First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	productDetail := dto.ProductDetailDTO{
		ID:               product.ID,
		Name:             product.Name,
		ShortDescription: product.ShortDescription,
		Description:      product.Description,
		Price:            product.Price,
		Stock:            product.Stock,
		CategoryID:       product.CategoryID,
		Category:         product.Category,
		GroupID:          product.GroupID,
		Group:            product.Group,
		SKU:              product.SKU,
		ImageURL:         product.ImageURL,
		Weight:           product.Weight,
		Dimensions:       product.Dimensions,
		Tags:             product.Tags,
		Status:           product.Status,
		Variants:         product.Variants,
		Addons:           product.Addons,
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
		LastSoldAt:       product.LastSoldAt,
	}

	ctx.JSON(http.StatusOK, productDetail)
}

func (c *ProductController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	tx := c.db.Begin()

	// Clear associations first
	if err := tx.Model(&models.Product{}).Where("id = ?", id).Association("Addons").Clear(); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear product associations"})
		return
	}

	// Delete variants
	if err := tx.Where("product_id = ?", id).Delete(&models.ProductVariant{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product variants"})
		return
	}

	// Delete the product
	if err := tx.Delete(&models.Product{}, id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// Helper functions
func groupExists(db *gorm.DB, id uint) bool {
	var exists bool
	db.Model(&models.ProductGroup{}).Select("count(*) > 0").Where("id = ?", id).Find(&exists)
	return exists
}

func categoryExists(db *gorm.DB, id uint) bool {
	var exists bool
	db.Model(&models.ProductCategory{}).Select("count(*) > 0").Where("id = ?", id).Find(&exists)
	return exists
}
