package dto

// dto/variant_dto.go
type CreateVariantDTO struct {
	ProductID  uint    `json:"productId" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	SKU        string  `json:"sku" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Stock      int     `json:"stock" binding:"required,gte=0"`
	ImageURL   string  `json:"imageUrl"`
	Attributes string  `json:"attributes"`
	IsDefault  bool    `json:"isDefault"`
	Weight     float64 `json:"weight"`
	Dimensions string  `json:"dimensions"`
	BarCode    string  `json:"barCode"`
}

type UpdateVariantDTO struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price" binding:"omitempty,gt=0"`
	Stock      int     `json:"stock" binding:"omitempty,gte=0"`
	ImageURL   string  `json:"imageUrl"`
	Attributes string  `json:"attributes"`
	IsDefault  *bool   `json:"isDefault"`
}
