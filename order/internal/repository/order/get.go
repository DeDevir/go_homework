package order

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/DeDevir/go_homework/order/internal/repository/converter"
	"github.com/google/uuid"
)

func (r *repository) Get(_ context.Context, uuid uuid.UUID) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	orderDto, ok := r.data[uuid.String()]
	if !ok {
		return nil, model.OrderNotFound
	}
	return converter.ParseOrderDtoToModel(orderDto), nil
}
