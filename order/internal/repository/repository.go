package repository

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, info model.Order) (*model.Order, error)
	Get(ctx context.Context, uuid uuid.UUID) (*model.Order, error)
	Cancel(ctx context.Context, uuid uuid.UUID) error
	Pay(ctx context.Context, uuid uuid.UUID, method model.PaymentMethod, transactionUUID uuid.UUID) error
}
