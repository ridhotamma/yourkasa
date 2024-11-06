package models

import (
	"gorm.io/gorm"
)

type ProductCategory struct {
	gorm.Model
	Name        string           `json:"name" gorm:"not null"`
	Slug        string           `json:"slug" gorm:"uniqueIndex;not null"`
	Description string           `json:"description"`
	ImageURL    string           `json:"imageUrl"`
	ParentID    *uint            `json:"parentId"` // For nested categories
	Parent      *ProductCategory `json:"parent" gorm:"foreignKey:ParentID"`
	Products    []Product        `json:"products"`
	Level       int              `json:"level" gorm:"not null"` // Depth level in category tree
	IsActive    bool             `json:"isActive" gorm:"default:true"`
	SortOrder   int              `json:"sortOrder" gorm:"default:0"`
}
