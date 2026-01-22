package payment

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuit struct {
	suite.Suite

	ctx context.Context

	service *service
}

func (s *ServiceSuit) SetupTest() {
	s.ctx = context.Background()

	s.service = NewService()
}

func (s *ServiceSuit) TearDownTest() {

}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuit))
}
