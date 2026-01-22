package order

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/DeDevir/go_homework/order/internal/repository/converter"
	"time"
)

func (r *repository) Create(_ context.Context, model model.Order) (*model.Order, error) {
	modelDto := converter.ParseOrderModelToDto(model)
	r.mu.Lock()
	defer r.mu.Unlock()
	modelDto.CreatedAt = time.Now()
	r.data[modelDto.UUID] = modelDto
	return converter.ParseOrderDtoToModel(modelDto), nil
}
