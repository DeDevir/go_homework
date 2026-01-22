package part

import (
	"context"
	"github.com/DeDevir/go_homework/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuit struct {
	suite.Suite

	ctx context.Context

	partRepository *mocks.PartRepository

	service *service
}

func (s *ServiceSuit) SetupTest() {
	s.ctx = context.Background()

	s.partRepository = mocks.NewPartRepository(s.T())
	s.service = NewService(s.partRepository)
}

func (s *ServiceSuit) TearDownTest() {

}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuit))
}
