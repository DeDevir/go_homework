package order

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

func (s *service) Pay(ctx context.Context, orderUuid uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	order, err := s.Get(ctx,orderUuid)
	if err != nil {
		return uuid.Nil,err
	}
	
	transactionUUID, err := s.paymentClient.PayOrder(ctx,order.UUID,order.UserUUID,method)
	if err != nil {
		grpcStat , ok := grpcStatus.FromError(err)
		
		if !ok {
			return uuid.Nil, model.InternalServerError
		}
		
		switch grpcStat.Code() {
		case codes.Internal:
			return uuid.Nil, model.InternalServerError
		default:
			return uuid.Nil, model.InternalServerError
		}
	}
	
	err = s.orderRepository.Pay(ctx, orderUuid, method,transactionUUID)
	if err != nil {
		return uuid.Nil, err
	}
	return transactionUUID, nil
}
