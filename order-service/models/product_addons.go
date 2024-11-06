package models

import (
	"gorm.io/gorm"
)

type ProductAddon struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	SKU         string    `json:"sku" gorm:"uniqueIndex;not null"`
	ImageURL    string    `json:"imageUrl"`
	Stock       int       `json:"stock" gorm:"not null"`
	IsRequired  bool      `json:"isRequired" gorm:"default:false"`
	MaxQuantity int       `json:"maxQuantity" gorm:"default:1"`
	Status      string    `json:"status" gorm:"type:varchar(20);default:'active'"`
	Products    []Product `json:"products" gorm:"many2many:product_addon_mappings;"`
}

type ProductAddonMapping struct {
	ProductID uint         `gorm:"primaryKey"`
	AddonID   uint         `gorm:"primaryKey"`
	Product   Product      `json:"-"`
	Addon     ProductAddon `json:"-"`
}
