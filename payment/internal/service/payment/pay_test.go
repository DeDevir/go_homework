package payment

import (
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuit) TestPaySuccess() {
	var (
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		methodPayment = "CARD"
	)

	transactionUUID, err := s.service.Pay(s.ctx, orderUUID, userUUID, methodPayment)

	s.Require().NoError(err)
	s.Require().NotEmpty(transactionUUID)
}
