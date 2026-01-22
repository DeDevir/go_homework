package order

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) Create(ctx context.Context, userUuid uuid.UUID, partsUuids uuid.UUIDs) (uuid.UUID, float64, error) {
	orderUuid := uuid.New()
	
	
	partFilter := model.PartFilter{
		UUIDS: partsUuids.Strings(),
	}
	parts, err := s.inventoryClient.ListParts(ctx,partFilter)
	if err != nil {
		return uuid.Nil, 0, err
	}
	var price float64
	for _,part := range parts {
		price += part.Price
	}
	
	order := model.Order{
		UUID:            orderUuid,
		UserUUID:        userUuid,
		PartUUIDs:       partsUuids,
		TotalPrice:      price,
		Status:          model.OrderStatusPENDINGPAYMENT,
	}
	
	orderSaved, err := s.orderRepository.Create(ctx,order)
	if err != nil {
		return uuid.Nil,0,err
	}
	return orderSaved.UUID, order.TotalPrice, nil
}