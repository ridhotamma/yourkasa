package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	Price       float64    `json:"price" gorm:"not null"`
	Stock       int        `json:"stock" gorm:"not null"`
	Category    string     `json:"category" gorm:"not null"`
	SKU         string     `json:"sku" gorm:"uniqueIndex;not null"`
	ImageURL    string     `json:"imageUrl"`
	LastSoldAt  *time.Time `json:"lastSoldAt"`
}
