package v1

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/converter"
	"github.com/DeDevir/go_homework/order/internal/model"
	orderV1 "github.com/DeDevir/go_homework/shared/pkg/openapi/order/v1"
	"github.com/go-faster/errors"
)

func (a *api) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	orderUuid := params.OrderUUID
	order, err := a.service.Get(ctx, orderUuid)
	if err != nil {
		if errors.Is(err, model.OrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order not found",
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}
	orderDto := converter.ParseOrderModelToOpenApi(order)
	return orderDto, nil
}
