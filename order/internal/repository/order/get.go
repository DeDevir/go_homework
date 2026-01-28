package order

import (
	"context"
	"errors"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/DeDevir/go_homework/order/internal/repository/converter"
	repoModel "github.com/DeDevir/go_homework/order/internal/repository/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Get(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	query := `
        SELECT
            uuid,
            user_uuid,
            part_uuids,
            total_price,
            transaction_uuid,
            payment_method,
            order_status,
            created_at,
            updated_at
        FROM orders
        WHERE uuid = $1;
    `

	var dto repoModel.Order

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&dto.UUID,
		&dto.UserUUID,
		&dto.PartUUIDs,
		&dto.TotalPrice,
		&dto.TransactionUUID,
		&dto.PaymentMethod,
		&dto.Status,
		&dto.CreatedAt,
		&dto.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.OrderNotFound
		}
		return nil, err
	}

	return converter.ParseOrderDtoToModel(&dto), nil
}
