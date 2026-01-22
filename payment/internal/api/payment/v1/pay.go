package v1

import (
	"context"
	paymentV1 "github.com/DeDevir/go_homework/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	uuid, err := a.service.Pay(ctx, req.OrderUuid, req.UserUuid, req.PaymentMethod.String())

	if err != nil {
		return nil, err
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: uuid,
	}, nil
}
