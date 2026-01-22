package order

import (
	"context"
	"github.com/google/uuid"
)

func (s *service) Cancel(ctx context.Context, uuid uuid.UUID) error {
	err := s.orderRepository.Cancel(ctx, uuid)
	if err != nil {
		return err
	}
	return nil
}