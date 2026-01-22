package part

import (
	"context"
	"github.com/DeDevir/go_homework/inventory/internal/model"
	"github.com/DeDevir/go_homework/inventory/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.data[uuid]
	if !ok {
		return nil, model.OrderNotFound
	}
	return converter.PartInfoToModel(part), nil
}
