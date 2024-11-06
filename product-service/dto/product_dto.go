package dto

import (
	"time"

	"github.com/ridhotamma/yourkasa/product-service/models"
)

type CreateProductDTO struct {
	Name             string  `json:"name" binding:"required"`
	ShortDescription string  `json:"shortDescription"`
	Description      string  `json:"description"`
	Price            float64 `json:"price" binding:"required,gt=0"`
	Stock            int     `json:"stock" binding:"required,gte=0"`
	CategoryID       uint    `json:"categoryId" binding:"required"`
	GroupID          *uint   `json:"groupId"`
	SKU              string  `json:"sku" binding:"required"`
	ImageURL         string  `json:"imageUrl"`
	Weight           float64 `json:"weight"`
	Dimensions       string  `json:"dimensions"`
	Tags             string  `json:"tags"`
	Status           string  `json:"status" binding:"required,oneof=active inactive discontinued"`
	AddonIDs         []uint  `json:"addonIds"`
}

type UpdateProductDTO struct {
	Name             string  `json:"name"`
	ShortDescription string  `json:"shortDescription"`
	Description      string  `json:"description"`
	Price            float64 `json:"price" binding:"omitempty,gt=0"`
	Stock            int     `json:"stock" binding:"omitempty,gte=0"`
	CategoryID       *uint   `json:"categoryId"`
	GroupID          *uint   `json:"groupId"`
	ImageURL         string  `json:"imageUrl"`
	Weight           float64 `json:"weight"`
	Dimensions       string  `json:"dimensions"`
	Tags             string  `json:"tags"`
	Status           string  `json:"status" binding:"omitempty,oneof=active inactive discontinued"`
	AddonIDs         []uint  `json:"addonIds"`
}

type ProductListDTO struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price"`
	Stock            int     `json:"stock"`
	CategoryID       uint    `json:"categoryId"`
	CategoryName     string  `json:"categoryName"`
	SKU              string  `json:"sku"`
	Status           string  `json:"status"`
	ImageURL         string  `json:"imageUrl"`
	VariantCount     int     `json:"variantCount"`
}

type ProductDetailDTO struct {
	ID               uint                    `json:"id"`
	Name             string                  `json:"name"`
	ShortDescription string                  `json:"shortDescription"`
	Description      string                  `json:"description"`
	Price            float64                 `json:"price"`
	Stock            int                     `json:"stock"`
	CategoryID       uint                    `json:"categoryId"`
	Category         models.ProductCategory  `json:"category"`
	GroupID          *uint                   `json:"groupId"`
	Group            *models.ProductGroup    `json:"group,omitempty"`
	SKU              string                  `json:"sku"`
	ImageURL         string                  `json:"imageUrl"`
	Weight           float64                 `json:"weight"`
	Dimensions       string                  `json:"dimensions"`
	Tags             string                  `json:"tags"`
	Status           string                  `json:"status"`
	Variants         []models.ProductVariant `json:"variants"`
	Addons           []models.ProductAddon   `json:"addons"`
	CreatedAt        time.Time               `json:"createdAt"`
	UpdatedAt        time.Time               `json:"updatedAt"`
	LastSoldAt       *time.Time              `json:"lastSoldAt"`
}

// Additional DTOs for nested responses
type ProductVariantDTO struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	ImageURL   string  `json:"imageUrl"`
	Attributes string  `json:"attributes"`
	IsDefault  bool    `json:"isDefault"`
	Status     string  `json:"status"`
	Weight     float64 `json:"weight"`
	Dimensions string  `json:"dimensions"`
	BarCode    string  `json:"barCode"`
}

type ProductAddonDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	IsRequired  bool    `json:"isRequired"`
	MaxQuantity int     `json:"maxQuantity"`
	Status      string  `json:"status"`
}

// Search and filter DTOs
type ProductSearchParams struct {
	CategoryID *uint   `form:"categoryId"`
	GroupID    *uint   `form:"groupId"`
	Status     string  `form:"status"`
	MinPrice   float64 `form:"minPrice"`
	MaxPrice   float64 `form:"maxPrice"`
	InStock    *bool   `form:"inStock"`
	Query      string  `form:"q"`
	SortBy     string  `form:"sortBy"`
	SortOrder  string  `form:"sortOrder"`
	Page       int     `form:"page,default=1"`
	PageSize   int     `form:"pageSize,default=20"`
}

type ProductListResponse struct {
	Products    []ProductListDTO `json:"products"`
	TotalCount  int64            `json:"totalCount"`
	PageCount   int              `json:"pageCount"`
	CurrentPage int              `json:"currentPage"`
	PageSize    int              `json:"pageSize"`
}
