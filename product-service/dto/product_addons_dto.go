package dto

type CreateAddonDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	SKU         string  `json:"sku" binding:"required"`
	ImageURL    string  `json:"imageUrl"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	IsRequired  bool    `json:"isRequired"`
	MaxQuantity int     `json:"maxQuantity"`
}

type UpdateAddonDTO struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
	IsRequired  *bool   `json:"isRequired"`
	MaxQuantity *int    `json:"maxQuantity"`
}
