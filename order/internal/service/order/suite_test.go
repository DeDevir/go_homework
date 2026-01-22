package order

import (
	"context"
	clientMock "github.com/DeDevir/go_homework/order/internal/client/grpc/mocks"
	repoMock "github.com/DeDevir/go_homework/order/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
	
	ctx context.Context
	
	orderRepository *repoMock.OrderRepository
	
	paymentClient *clientMock.PaymentClient
	inventoryClient *clientMock.InventoryClient
	
	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	
	s.orderRepository = repoMock.NewOrderRepository(s.T())
	s.paymentClient = clientMock.NewPaymentClient(s.T())
	s.inventoryClient = clientMock.NewInventoryClient(s.T())
	
	s.service = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}