package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/product-service/dto"
	"github.com/ridhotamma/yourkasa/product-service/models"
	"gorm.io/gorm"
)

type CategoryController struct {
	db *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{db: db}
}

func (c *CategoryController) Create(ctx *gin.Context) {
	var input dto.CreateCategoryDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if slug exists
	var existingCategory models.ProductCategory
	if err := c.db.Where("slug = ?", input.Slug).First(&existingCategory).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Slug already exists"})
		return
	}

	category := models.ProductCategory{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		ParentID:    input.ParentID,
		Level:       input.Level,
		IsActive:    input.IsActive,
		SortOrder:   input.SortOrder,
	}

	if err := c.db.Create(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "id": category.ID})
}

func (c *CategoryController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var input dto.UpdateCategoryDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.ProductCategory
	if err := c.db.First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.ImageURL != "" {
		updates["image_url"] = input.ImageURL
	}
	if input.ParentID != nil {
		updates["parent_id"] = input.ParentID
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	if input.SortOrder != nil {
		updates["sort_order"] = *input.SortOrder
	}

	if err := c.db.Model(&category).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

func (c *CategoryController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	// Check for products in this category
	var productsCount int64
	if err := c.db.Model(&models.Product{}).Where("category_id = ?", id).Count(&productsCount).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check category usage"})
		return
	}

	if productsCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with associated products"})
		return
	}

	if err := c.db.Delete(&models.ProductCategory{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func (c *CategoryController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var category models.ProductCategory
	if err := c.db.Preload("Parent").First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (c *CategoryController) List(ctx *gin.Context) {
	var categories []models.ProductCategory
	if err := c.db.Preload("Parent").Find(&categories).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}
