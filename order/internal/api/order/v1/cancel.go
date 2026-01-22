package v1

import (
	"context"
	"errors"
	"github.com/DeDevir/go_homework/order/internal/model"
	orderV1 "github.com/DeDevir/go_homework/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	orderUuid := params.OrderUUID
	err := a.service.Cancel(ctx, orderUuid)
	if err != nil {
		if errors.Is(err, model.OrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order not found",
			}, nil
		}
		if errors.Is(err, model.OrderCannotBeCanceled) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Order has status paid and cannot be canceled",
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}
	return &orderV1.CancelOrderNoContent{}, nil
}
