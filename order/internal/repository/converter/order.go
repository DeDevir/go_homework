package converter

import (
	"github.com/DeDevir/go_homework/order/internal/model"
	repoModel "github.com/DeDevir/go_homework/order/internal/repository/model"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"log"
)

func ParseOrderModelToDto(order model.Order) *repoModel.Order {
	var PartUUIDsArray uuid.UUIDs = order.PartUUIDs
	return &repoModel.Order{
		UUID:            order.UUID.String(),
		UserUUID:        order.UserUUID.String(),
		PartUUIDs:       PartUUIDsArray.Strings(),
		TotalPrice:      order.TotalPrice,
		TransactionUUID: parseTransactionUUIDModelToDto(order.TransactionUUID),
		PaymentMethod:   ParsePaymentMethodModelToDto(order.PaymentMethod),
		Status:          parsOrderStatusModelToDto(order.Status),
	}
}

func parseTransactionUUIDModelToDto(trxUUID *uuid.UUID) *string {
	if trxUUID == nil {
		return nil
	}
	return lo.ToPtr(trxUUID.String())
}

func ParsePaymentMethodModelToDto(method *model.PaymentMethod) *repoModel.PaymentMethod {
	if method == nil {
		return lo.ToPtr(repoModel.PaymentMethodUnknown)
	}
	switch *method {
	case model.PaymentMethodCARD:
		return lo.ToPtr(repoModel.PaymentMethodCard)
	case model.PaymentMethodSBP:
		return lo.ToPtr(repoModel.PaymentMethodSBP)
	case model.PaymentMethodINVESTORMONEY:
		return lo.ToPtr(repoModel.PaymentMethodINVESTORMONEY)
	case model.PaymentMethodCREDITCARD:
		return lo.ToPtr(repoModel.PaymentMethodCREDITCARD)
	default:
		return lo.ToPtr(repoModel.PaymentMethodUnknown)
	}
}

func parsOrderStatusModelToDto(status model.OrderStatus) repoModel.OrderStatus {
	switch status {
	case model.OrderStatusPAID:
		return repoModel.OrderPaid
	case model.OrderStatusCANCELLED:
		return repoModel.OrderCanceled
	default:
		return repoModel.OrderPending
	}
}

func ParseOrderDtoToModel(order *repoModel.Order) *model.Order {
	partsUUIDs := parsePartsUUID(order.PartUUIDs)
	transactionUUID := parseTransactionUUIDDtoToModel(order.TransactionUUID)
	return &model.Order{
		UUID:            uuid.MustParse(order.UUID),
		UserUUID:        uuid.MustParse(order.UserUUID),
		PartUUIDs:       partsUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   parsePaymentMethodDtoToModel(order.PaymentMethod),
		Status:          parseStatusDtoToModel(order.Status),
	}
}

func parsePartsUUID(partsUUIDString []string) []uuid.UUID {
	parts := make([]uuid.UUID, 0, len(partsUUIDString))

	for _, v := range partsUUIDString {
		parts = append(parts, uuid.MustParse(v))
	}

	return parts
}
func parseTransactionUUIDDtoToModel(stringUuid *string) *uuid.UUID {
	if stringUuid == nil {
		return nil
	}
	transactionUuid := uuid.MustParse(*stringUuid)
	return &transactionUuid
}
func parsePaymentMethodDtoToModel(method *repoModel.PaymentMethod) *model.PaymentMethod {
	if method == nil {
		return nil
	}
	switch *method {
	case repoModel.PaymentMethodCard:
		return lo.ToPtr(model.PaymentMethodCARD)
	case repoModel.PaymentMethodCREDITCARD:
		return lo.ToPtr(model.PaymentMethodCREDITCARD)
	case repoModel.PaymentMethodINVESTORMONEY:
		return lo.ToPtr(model.PaymentMethodINVESTORMONEY)
	case repoModel.PaymentMethodSBP:
		return lo.ToPtr(model.PaymentMethodSBP)
	default:
		return lo.ToPtr(model.PaymentMethodUNKNOW)
	}
}
func parseStatusDtoToModel(status repoModel.OrderStatus) model.OrderStatus {
	switch status {
	case repoModel.OrderPaid:
		return model.OrderStatusPAID
	case repoModel.OrderCanceled:
		return model.OrderStatusCANCELLED
	default:
		log.Printf("default case (%v)", status)
		return model.OrderStatusPENDINGPAYMENT
	}
}
