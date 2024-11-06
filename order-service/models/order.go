package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderNumber     string      `json:"orderNumber" gorm:"uniqueIndex;not null"`
	CustomerID      string      `json:"customerId" gorm:"not null;index"`
	Status          string      `json:"status" gorm:"type:varchar(50);default:'pending'"`
	TotalAmount     float64     `json:"totalAmount" gorm:"not null"`
	SubtotalAmount  float64     `json:"subtotalAmount" gorm:"not null"`
	DiscountAmount  float64     `json:"discountAmount"`
	TaxAmount       float64     `json:"taxAmount"`
	ShippingAmount  float64     `json:"shippingAmount"`
	ShippingAddress string      `json:"shippingAddress" gorm:"type:text"`
	BillingAddress  string      `json:"billingAddress" gorm:"type:text"`
	PaymentMethod   string      `json:"paymentMethod"`
	PaymentStatus   string      `json:"paymentStatus" gorm:"type:varchar(50);default:'unpaid'"`
	Notes           string      `json:"notes"`
	OrderItems      []OrderItem `json:"orderItems"`
	PaidAt          *time.Time  `json:"paidAt"`
	CanceledAt      *time.Time  `json:"canceledAt"`
	CompletedAt     *time.Time  `json:"completedAt"`
}
