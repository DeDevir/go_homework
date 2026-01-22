package order

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	repoModel "github.com/DeDevir/go_homework/order/internal/repository/model"
	"github.com/google/uuid"
)

func (r *repository) Cancel(_ context.Context, uuid uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[uuid.String()]
	if !ok {
		return model.OrderNotFound
	}
	if order.Status == repoModel.OrderPaid {
		return model.OrderCannotBeCanceled
	}
	order.Status = repoModel.OrderCanceled
	return nil
}
