package order

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/DeDevir/go_homework/order/internal/repository/converter"
	repoModel "github.com/DeDevir/go_homework/order/internal/repository/model"
	"github.com/samber/lo"

	"github.com/google/uuid"
)

func (r *repository) Pay(_ context.Context, orderUuid uuid.UUID, method model.PaymentMethod, transactionUUID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.data[orderUuid.String()]
	if !ok {
		return model.OrderNotFound
	}
	if order.Status != repoModel.OrderPending {
		return model.OrderCannotBeCanceled
	}
	order.Status = repoModel.OrderPaid
	order.TransactionUUID = lo.ToPtr(transactionUUID.String())
	order.PaymentMethod = converter.ParsePaymentMethodModelToDto(&method)

	return nil
}
