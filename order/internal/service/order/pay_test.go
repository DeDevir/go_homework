package order

import (
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"log"
)

func (s *ServiceSuite) TestPayTestSuccess() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		userUUID = uuid.MustParse(gofakeit.UUID())
		order = &model.Order{
			UUID:            orderUUID ,
			UserUUID:        userUUID,
			PartUUIDs:       []uuid.UUID {
				uuid.New(),
				uuid.New(),
			},
			TotalPrice:      gofakeit.Price(100,10000),
			TransactionUUID: nil,
			PaymentMethod:   nil,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}
		transactionUUID = uuid.MustParse(gofakeit.UUID())
		errFindOrder error = nil
		errPayOrder error = nil
		methodPay = model.PaymentMethodSBP
	)
	
	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order,errFindOrder)
	s.paymentClient.On("PayOrder", s.ctx,order.UUID,order.UserUUID,methodPay).Return(transactionUUID,errPayOrder)
	s.orderRepository.On("Pay", s.ctx, order.UUID, methodPay, transactionUUID).Return(nil)
	
	trxUUID, err := s.service.Pay(s.ctx,order.UUID, methodPay)
	
	s.Require().NoError(err)
	s.Require().Equal(transactionUUID, trxUUID)
	
	order2, _ := s.service.Get(s.ctx, orderUUID)
	log.Printf("%#v", order2)
}

func (s *ServiceSuite) TestPayTestErrorClient() {

}