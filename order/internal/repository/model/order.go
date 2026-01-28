package model

import "time"

type Order struct {
	UUID            string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type OrderStatus string

const (
	OrderPaid     OrderStatus = "PAID"
	OrderCanceled OrderStatus = "CANCELLED"
	OrderPending  OrderStatus = "PENDING_PAYMENT"
)

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
)
