package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/product-service/dto"
	"github.com/ridhotamma/yourkasa/product-service/models"
	"gorm.io/gorm"
)

type GroupController struct {
	db *gorm.DB
}

func NewGroupController(db *gorm.DB) *GroupController {
	return &GroupController{db: db}
}

func (c *GroupController) Create(ctx *gin.Context) {
	var input dto.CreateGroupDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group := models.ProductGroup{
		Name:        input.Name,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		IsActive:    input.IsActive,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		GroupType:   input.GroupType,
		SortOrder:   input.SortOrder,
	}

	if err := c.db.Create(&group).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Group created successfully", "id": group.ID})
}

func (c *GroupController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var input dto.UpdateGroupDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var group models.ProductGroup
	if err := c.db.First(&group, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
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
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	if input.StartDate != nil {
		updates["start_date"] = input.StartDate
	}
	if input.EndDate != nil {
		updates["end_date"] = input.EndDate
	}
	if input.GroupType != "" {
		updates["group_type"] = input.GroupType
	}
	if input.SortOrder != nil {
		updates["sort_order"] = *input.SortOrder
	}

	if err := c.db.Model(&group).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Group updated successfully"})
}

func (c *GroupController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	// Check for products in this group
	var productsCount int64
	if err := c.db.Model(&models.Product{}).Where("group_id = ?", id).Count(&productsCount).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check group usage"})
		return
	}

	if productsCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete group with associated products"})
		return
	}

	if err := c.db.Delete(&models.ProductGroup{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}

func (c *GroupController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var group models.ProductGroup
	if err := c.db.Preload("Products").First(&group, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	ctx.JSON(http.StatusOK, group)
}

func (c *GroupController) List(ctx *gin.Context) {
	var groups []models.ProductGroup
	query := c.db.Model(&models.ProductGroup{})

	// Filter by active status if specified
	if activeStr := ctx.Query("active"); activeStr != "" {
		if active := activeStr == "true"; active {
			query = query.Where("is_active = ?", true)
		}
	}

	// Filter by group type if specified
	if groupType := ctx.Query("type"); groupType != "" {
		query = query.Where("group_type = ?", groupType)
	}

	// Filter by date range
	if err := query.Preload("Products").Order("sort_order asc").Find(&groups).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}

	ctx.JSON(http.StatusOK, groups)
}

// Additional helper methods for GroupController

func (c *GroupController) AddProductToGroup(ctx *gin.Context) {
	groupID := ctx.Param("id")
	var input struct {
		ProductID uint `json:"productId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify group exists
	var group models.ProductGroup
	if err := c.db.First(&group, groupID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Verify product exists
	var product models.Product
	if err := c.db.First(&product, input.ProductID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Update product's group ID
	if err := c.db.Model(&product).Update("group_id", groupID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to group"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product added to group successfully"})
}

func (c *GroupController) RemoveProductFromGroup(ctx *gin.Context) {
	groupID := ctx.Param("id")
	productID := ctx.Param("productId")

	// Verify group exists
	var group models.ProductGroup
	if err := c.db.First(&group, groupID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Remove product from group
	if err := c.db.Model(&models.Product{}).Where("id = ? AND group_id = ?", productID, groupID).
		Update("group_id", nil).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove product from group"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product removed from group successfully"})
}
