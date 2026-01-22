package part

import (
	"context"
	"github.com/DeDevir/go_homework/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.partRepository.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return part, nil
}
