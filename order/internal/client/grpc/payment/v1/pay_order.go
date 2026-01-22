package v1

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/client/converter"
	"github.com/DeDevir/go_homework/order/internal/model"
	paymentV1 "github.com/DeDevir/go_homework/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"log"
)

func (c *client) PayOrder(ctx context.Context, orderUuid uuid.UUID, userUuid uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	resp, err := c.generatedClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderUuid.String(),
		UserUuid:      userUuid.String(),
		PaymentMethod: converter.MethodPaymentModelToProto(method),
	})

	if err != nil {
		log.Printf("Client error payorder err:%v\n", err)
		return uuid.UUID{}, err
	}
	transactionUUID := resp.TransactionUuid
	return uuid.MustParse(transactionUUID), nil
}
