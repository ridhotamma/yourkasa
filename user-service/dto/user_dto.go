package dto

import "time"

type CreateUserDTO struct {
	FirstName         string `json:"firstName" binding:"required"`
	LastName          string `json:"lastName" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	Password          string `json:"password" binding:"required,min=6"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
	Role              string `json:"role" binding:"required,oneof=admin cashier owner"`
}

type UpdateUserDTO struct {
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
	Role              string `json:"role" binding:"omitempty,oneof=admin cashier owner"`
}

type UserListDTO struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type UserDetailDTO struct {
	ID                uint       `json:"id"`
	FirstName         string     `json:"firstName"`
	LastName          string     `json:"lastName"`
	Email             string     `json:"email"`
	ProfilePictureUrl string     `json:"profilePictureUrl"`
	Role              string     `json:"role"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	LastLoggedIn      *time.Time `json:"lastLoggedIn"`
}
