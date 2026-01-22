package part

import (
	"context"
	"github.com/DeDevir/go_homework/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filters *model.PartsFilter) ([]*model.Part, error) {
	parts, err := s.partRepository.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
