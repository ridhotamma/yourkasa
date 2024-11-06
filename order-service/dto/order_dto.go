package dto

type CreateOrderDTO struct {
	ShippingAddress string `json:"shippingAddress" binding:"required"`
	BillingAddress  string `json:"billingAddress" binding:"required"`
	PaymentMethod   string `json:"paymentMethod" binding:"required"`
	Notes           string `json:"notes"`
}
