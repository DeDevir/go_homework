package order

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/DeDevir/go_homework/order/internal/repository/converter"
	"log"
)

func (r *repository) Create(ctx context.Context, model model.Order) (*model.Order, error) {
	modelDto := converter.ParseOrderModelToDto(model)
	query := `
        INSERT INTO orders (user_uuid, part_uuids, total_price, order_status)
        VALUES ($1, $2, $3, $4)
        RETURNING uuid, user_uuid, part_uuids, total_price, transaction_uuid, payment_method, order_status, created_at, updated_at;
    `

	err := r.pool.QueryRow(ctx, query,
		model.UserUUID,
		model.PartUUIDs,
		model.TotalPrice,
		model.Status, // должно быть OrderPending
	).Scan(
		&modelDto.UUID,
		&modelDto.UserUUID,
		&modelDto.PartUUIDs,
		&modelDto.TotalPrice,
		&modelDto.TransactionUUID,
		&modelDto.PaymentMethod,
		&modelDto.Status,
		&modelDto.CreatedAt,
		&modelDto.UpdatedAt,
	)

	if err != nil {
		log.Printf("Failed to insert order: %v\n", err)
		return nil, err
	}

	return converter.ParseOrderDtoToModel(modelDto), nil
}
