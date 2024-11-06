package models

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	OrderID     uint    `json:"orderId" gorm:"not null"`
	ProductID   uint    `json:"productId" gorm:"not null"`
	VariantID   *uint   `json:"variantId"`
	ProductName string  `json:"productName" gorm:"not null"`
	VariantName string  `json:"variantName"`
	Quantity    int     `json:"quantity" gorm:"not null"`
	Price       float64 `json:"price" gorm:"not null"`
	SubTotal    float64 `json:"subTotal" gorm:"not null"`
	AddonsData  string  `json:"addonsData" gorm:"type:jsonb"`
	Order       Order   `json:"-"`
}
