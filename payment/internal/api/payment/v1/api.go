package v1

import (
	"github.com/DeDevir/go_homework/payment/internal/service"
	paymentV1 "github.com/DeDevir/go_homework/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	service service.PaymentService
}

func NewApi(paymentService service.PaymentService) *api {
	return &api{
		service: paymentService,
	}
}
