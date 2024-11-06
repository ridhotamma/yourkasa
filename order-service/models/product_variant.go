package models

import (
	"gorm.io/gorm"
)

type ProductVariant struct {
	gorm.Model
	ProductID  uint    `json:"productId" gorm:"not null"`
	Name       string  `json:"name" gorm:"not null"`
	SKU        string  `json:"sku" gorm:"uniqueIndex;not null"`
	Price      float64 `json:"price" gorm:"not null"`
	Stock      int     `json:"stock" gorm:"not null"`
	ImageURL   string  `json:"imageUrl"`
	Attributes string  `json:"attributes" gorm:"type:jsonb"` // JSON string storing variant attributes (color, size, etc.)
	IsDefault  bool    `json:"isDefault" gorm:"default:false"`
	Status     string  `json:"status" gorm:"type:varchar(20);default:'active'"`
	Weight     float64 `json:"weight" gorm:"type:decimal(10,2)"`
	Dimensions string  `json:"dimensions"`
	BarCode    string  `json:"barCode"`
	Product    Product `json:"-"`
}
