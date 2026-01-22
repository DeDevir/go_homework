package service

import "context"

type PaymentService interface {
	Pay(ctx context.Context, orderUUID string, userUUID string, methodPayment string) (string, error)
}
