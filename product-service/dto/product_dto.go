package dto

import "time"

type CreateProductDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	Category    string  `json:"category" binding:"required"`
	SKU         string  `json:"sku" binding:"required"`
	ImageURL    string  `json:"imageUrl"`
}

type UpdateProductDTO struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"imageUrl"`
}

type ProductListDTO struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	Category string  `json:"category"`
	SKU      string  `json:"sku"`
}

type ProductDetailDTO struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	Category    string     `json:"category"`
	SKU         string     `json:"sku"`
	ImageURL    string     `json:"imageUrl"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastSoldAt  *time.Time `json:"lastSoldAt"`
}
