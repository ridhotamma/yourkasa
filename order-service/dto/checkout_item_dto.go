package dto

type CreateCheckoutItemDTO struct {
	ProductID  uint               `json:"productId" binding:"required"`
	VariantID  *uint              `json:"variantId"`
	Quantity   int                `json:"quantity" binding:"required,gt=0"`
	AddonsData []CheckoutAddonDTO `json:"addonsData"`
	Notes      string             `json:"notes"`
}

type CheckoutAddonDTO struct {
	AddonID  uint `json:"addonId" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,gt=0"`
}

type UpdateCheckoutItemDTO struct {
	Quantity   int                `json:"quantity" binding:"omitempty,gt=0"`
	AddonsData []CheckoutAddonDTO `json:"addonsData"`
	Notes      string             `json:"notes"`
	IsSelected *bool              `json:"isSelected"`
}
