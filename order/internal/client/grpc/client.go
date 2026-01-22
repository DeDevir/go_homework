package grpc

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/google/uuid"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartFilter) ([]*model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUuid uuid.UUID, userUuid uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
}
