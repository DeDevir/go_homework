package converter // internal.converter

import (
	"github.com/DeDevir/go_homework/order/internal/model"
	orderV1 "github.com/DeDevir/go_homework/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func ParseOrderOpenApiToModel(dto orderV1.OrderDto) *model.Order {
	return &model.Order{
		UUID:            dto.OrderUUID,
		UserUUID:        dto.UserUUID,
		PartUUIDs:       dto.PartUuids,
		TotalPrice:      dto.TotalPrice,
		TransactionUUID: &dto.TransactionUUID.Value,
		PaymentMethod:   nil,
		Status:          parseOrderStatusOpenApiToModel(dto.Status),
	}
}

func ParseOrderModelToOpenApi(model *model.Order) *orderV1.OrderDto {
	return &orderV1.OrderDto{
		OrderUUID:       model.UUID,
		UserUUID:        model.UserUUID,
		PartUuids:       model.PartUUIDs,
		TotalPrice:      model.TotalPrice,
		TransactionUUID: parseTransactionUUIDModelToOpenApi(model.TransactionUUID),
		PaymentMethod:   orderV1.OptString{},
		Status:          parseOrderStatusModelToOpenApi(model.Status),
	}
}

func parseTransactionUUIDModelToOpenApi(transactionUuid *uuid.UUID) orderV1.OptUUID {
	if transactionUuid == nil {
		return orderV1.OptUUID{Set: false}
	}
	return orderV1.OptUUID{
		Value: *transactionUuid,
		Set:   true,
	}
}

func parseOrderStatusOpenApiToModel(status orderV1.OrderStatus) model.OrderStatus {
	switch status {
	case orderV1.OrderStatusCANCELLED:
		return model.OrderStatusCANCELLED
	case orderV1.OrderStatusPAID:
		return model.OrderStatusPAID
	default:
		return model.OrderStatusPENDINGPAYMENT
	}
}

func parseOrderStatusModelToOpenApi(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusPAID:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCANCELLED:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

func ParsePaymentMethodOpenApiToModel(method orderV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	case orderV1.PaymentMethodPAYMENTMETHODCARD:
		return model.PaymentMethodCARD
	case orderV1.PaymentMethodPAYMENTMETHODSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodPAYMENTMETHODCREDITCARD:
		return model.PaymentMethodCREDITCARD
	default:
		return model.PaymentMethodUNKNOW
	}
}
