package v1

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/converter"
	"github.com/DeDevir/go_homework/order/internal/model"
	orderV1 "github.com/DeDevir/go_homework/shared/pkg/openapi/order/v1"
	"github.com/go-faster/errors"
	"log"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	transactionUuid, err := a.service.Pay(ctx, params.OrderUUID, converter.ParsePaymentMethodOpenApiToModel(req.PaymentMethod))
	if err != nil {
		if errors.Is(err, model.OrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order not found",
			}, nil
		}
		log.Printf("Error pay: %v\n", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}
	return &orderV1.PayOrderResponse{TransactionUUID: transactionUuid}, nil
}
