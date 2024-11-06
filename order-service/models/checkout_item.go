package models

import (
	"gorm.io/gorm"
)

type CheckoutItem struct {
	gorm.Model
	CustomerID string          `json:"customerId" gorm:"not null"`
	ProductID  uint            `json:"productId" gorm:"not null"`
	VariantID  *uint           `json:"variantId"`
	Quantity   int             `json:"quantity" gorm:"not null"`
	Price      float64         `json:"price" gorm:"not null"`
	Product    Product         `json:"product"`
	Variant    *ProductVariant `json:"variant"`
	AddonsData string          `json:"addonsData" gorm:"type:jsonb"` // Stores selected addons and their quantities
	Notes      string          `json:"notes"`
	IsSelected bool            `json:"isSelected" gorm:"default:true"`
}
