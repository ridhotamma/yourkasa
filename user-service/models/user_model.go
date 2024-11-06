package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleCashier Role = "cashier"
	RoleOwner   Role = "owner"
)

type User struct {
	gorm.Model
	FirstName         string     `json:"firstName" gorm:"not null"`
	LastName          string     `json:"lastName" gorm:"not null"`
	Email             string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash      string     `json:"-" gorm:"not null"`
	ProfilePictureUrl string     `json:"profilePictureUrl"`
	Role              Role       `json:"role" gorm:"type:varchar(10);not null"`
	LastLoggedIn      *time.Time `json:"lastLoggedIn"`
}
