package v1

import (
	"github.com/DeDevir/go_homework/order/internal/service"
	orderV1 "github.com/DeDevir/go_homework/shared/pkg/openapi/order/v1"
)

type api struct {
	service service.OrderService
	orderV1.UnimplementedHandler
}

func NewApi(service service.OrderService) *api {
	return &api{
		service: service,
	}
}
