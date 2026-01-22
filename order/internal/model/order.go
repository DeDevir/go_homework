package model // internal.model.order

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodUNKNOW        PaymentMethod = "PAYMENT_METHOD_UNKNOWN_UNSPECIFIED"
	PaymentMethodCARD          PaymentMethod = "PAYMENT_METHOD_CARD"
	PaymentMethodSBP           PaymentMethod = "PAYMENT_METHOD_SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "PAYMENT_METHOD_CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "PAYMENT_METHOD_INVESTOR_MONEY"
)

type Order struct {
	UUID            uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []uuid.UUID
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
}
