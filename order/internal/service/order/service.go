package order

import (
	"github.com/DeDevir/go_homework/order/internal/client/grpc"
	"github.com/DeDevir/go_homework/order/internal/repository"
	def "github.com/DeDevir/go_homework/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository
	
	inventoryClient grpc.InventoryClient
	paymentClient grpc.PaymentClient
}

func NewService(repository repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient) *service {
	return &service{
		orderRepository: repository,
		inventoryClient: inventoryClient,
		paymentClient: paymentClient,
	}
}