package converter

import (
	"github.com/DeDevir/go_homework/order/internal/model"
	paymentV1 "github.com/DeDevir/go_homework/shared/pkg/proto/payment/v1"
)

func MethodPaymentModelToProto(method model.PaymentMethod) paymentV1.PaymentMethod {
	switch method {
	case model.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case model.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
	}
}
