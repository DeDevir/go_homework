package order

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) Get(ctx context.Context, uuid uuid.UUID) (*model.Order, error) {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return nil,err
	}
	return order,err
}