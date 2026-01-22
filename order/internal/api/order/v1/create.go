package v1

import (
	"context"
	"errors"
	"github.com/DeDevir/go_homework/order/internal/model"
	orderV1 "github.com/DeDevir/go_homework/shared/pkg/openapi/order/v1"
)

func (a *api) NewOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.NewOrderRes, error) {
	userUuid := req.UserUUID
	partsUuid := req.PartUuids
	orderUuid, price, err := a.service.Create(ctx, userUuid, partsUuid)
	if err != nil {
		if errors.Is(err, model.SearchablePartsNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Requested parts not found",
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}
	return &orderV1.CreateOrderResponse{
		UUID:       orderV1.NewOptUUID(orderUuid),
		TotalPrice: price,
	}, nil
}
