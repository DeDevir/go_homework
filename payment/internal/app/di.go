package app

import (
	"context"
	paymentV1Api "github.com/DeDevir/go_homework/payment/internal/api/payment/v1"
	"github.com/DeDevir/go_homework/payment/internal/service"
	paymentService "github.com/DeDevir/go_homework/payment/internal/service/payment"
	paymentV1 "github.com/DeDevir/go_homework/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API paymentV1.PaymentServiceServer

	paymentService service.PaymentService
}

func NewDIContainer() *diContainer { return &diContainer{} }

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentV1Api.NewApi(d.PaymentService(ctx))
	}
	return d.paymentV1API
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}
	return d.paymentService
}
