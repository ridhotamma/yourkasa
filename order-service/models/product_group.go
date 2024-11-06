package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductGroup struct {
	gorm.Model
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	ImageURL    string     `json:"imageUrl"`
	Products    []Product  `json:"products"`
	IsActive    bool       `json:"isActive" gorm:"default:true"`
	StartDate   *time.Time `json:"startDate"`                         // For temporary/seasonal groups
	EndDate     *time.Time `json:"endDate"`                           // For temporary/seasonal groups
	GroupType   string     `json:"groupType" gorm:"type:varchar(50)"` // bundle, collection, seasonal, etc.
	SortOrder   int        `json:"sortOrder" gorm:"default:0"`
}
