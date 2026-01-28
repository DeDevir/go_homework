package part

import (
	"context"
	"github.com/DeDevir/go_homework/inventory/internal/model"
	"log"
)

func (s *service) Get(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.partRepository.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	log.Printf("finded part in service %v", part)

	return part, nil
}
