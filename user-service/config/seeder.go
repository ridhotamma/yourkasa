package config

import (
	"log"

	"github.com/ridhotamma/yourkasa/user-service/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type defaultUser struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Role      models.Role
}

func SeedUsers(db *gorm.DB) error {
	defaultUsers := []defaultUser{
		{
			FirstName: "Admin",
			LastName:  "User",
			Email:     "admin@yourkasa.com",
			Password:  "admin123",
			Role:      models.RoleAdmin,
		},
		{
			FirstName: "Cashier",
			LastName:  "User",
			Email:     "cashier@yourkasa.com",
			Password:  "cashier123",
			Role:      models.RoleCashier,
		},
		{
			FirstName: "Owner",
			LastName:  "User",
			Email:     "owner@yourkasa.com",
			Password:  "owner123",
			Role:      models.RoleOwner,
		},
	}

	for _, user := range defaultUsers {
		// Check if user already exists
		var existingUser models.User
		result := db.Where("email = ?", user.Email).First(&existingUser)
		if result.Error == nil {
			log.Printf("User with email %s already exists, skipping...", user.Email)
			continue
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Create new user
		newUser := models.User{
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        user.Email,
			PasswordHash: string(hashedPassword),
			Role:         user.Role,
			LastLoggedIn: nil,
		}

		if err := db.Create(&newUser).Error; err != nil {
			return err
		}

		log.Printf("Created default user with email: %s and role: %s", user.Email, user.Role)
	}

	return nil
}
