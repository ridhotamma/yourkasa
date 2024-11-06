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

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		Category:    input.Category,
		SKU:         input.SKU,
		ImageURL:    input.ImageURL,
	}

	if err := c.db.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "id": product.ID})
}

func (c *ProductController) Update(ctx *gin.Context) {
	var input dto.UpdateProductDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	var product models.Product
	if err := c.db.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
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
	if input.Category != "" {
		updates["category"] = input.Category
	}
	if input.ImageURL != "" {
		updates["image_url"] = input.ImageURL
	}

	if err := c.db.Model(&product).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func (c *ProductController) List(ctx *gin.Context) {
	var products []models.Product
	if err := c.db.Find(&products).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	var productList []dto.ProductListDTO
	for _, product := range products {
		productList = append(productList, dto.ProductListDTO{
			ID:       product.ID,
			Name:     product.Name,
			Price:    product.Price,
			Stock:    product.Stock,
			Category: product.Category,
			SKU:      product.SKU,
		})
	}

	ctx.JSON(http.StatusOK, productList)
}

func (c *ProductController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.Product
	if err := c.db.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	productDetail := dto.ProductDetailDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		SKU:         product.SKU,
		ImageURL:    product.ImageURL,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		LastSoldAt:  product.LastSoldAt,
	}

	ctx.JSON(http.StatusOK, productDetail)
}

func (c *ProductController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.db.Delete(&models.Product{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
