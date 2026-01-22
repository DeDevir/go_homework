package service

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/google/uuid"
)

type OrderService interface {
	Get(ctx context.Context, uuid uuid.UUID) (*model.Order, error)
	Create(ctx context.Context, userUuid uuid.UUID, partUuids uuid.UUIDs) (uuid uuid.UUID, totalPrice float64, err error)
	Pay(ctx context.Context, uuid uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
	Cancel(ctx context.Context, uuid uuid.UUID) error
}
