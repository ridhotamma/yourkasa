package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name             string           `json:"name" gorm:"not null"`
	Description      string           `json:"description" gorm:"type:text"` // Extended description
	ShortDescription string           `json:"shortDescription" gorm:"type:varchar(255)"`
	Price            float64          `json:"price" gorm:"not null"`
	Stock            int              `json:"stock" gorm:"not null"`
	SKU              string           `json:"sku" gorm:"uniqueIndex;not null"`
	ImageURL         string           `json:"imageUrl"`
	LastSoldAt       *time.Time       `json:"lastSoldAt"`
	CategoryID       uint             `json:"categoryId" gorm:"not null"`
	Category         ProductCategory  `json:"category"`
	GroupID          *uint            `json:"groupId"`
	Group            *ProductGroup    `json:"group"`
	Variants         []ProductVariant `json:"variants"`
	Addons           []ProductAddon   `json:"addons" gorm:"many2many:product_addon_mappings;"`
	Status           string           `json:"status" gorm:"type:varchar(20);default:'active'"` // active, inactive, discontinued
	Weight           float64          `json:"weight" gorm:"type:decimal(10,2)"`                // in kg
	Dimensions       string           `json:"dimensions"`                                      // JSON string storing length, width, height
	Tags             string           `json:"tags"`                                            // Comma-separated tags
}
