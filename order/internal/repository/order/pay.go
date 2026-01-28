package order

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/DeDevir/go_homework/order/internal/repository/converter"
	repoModel "github.com/DeDevir/go_homework/order/internal/repository/model"

	"github.com/google/uuid"
)

func (r *repository) Pay(
	ctx context.Context,
	orderUUID uuid.UUID,
	method model.PaymentMethod,
	transactionUUID uuid.UUID,
) error {

	query := `
        UPDATE orders
        SET
            order_status = $2,
            transaction_uuid = $3,
            payment_method = $4
        WHERE
            uuid = $1
            AND order_status = $5;
    `

	commandTag, err := r.pool.Exec(
		ctx,
		query,
		orderUUID,
		repoModel.OrderPaid,
		transactionUUID,
		converter.ParsePaymentMethodModelToDto(&method),
		repoModel.OrderPending,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		// либо заказа нет, либо он не в PENDING_PAYMENT
		return model.OrderAlreadyPaid
	}

	return nil
}
