package order

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/model"
	"github.com/google/uuid"
	"log"
)

func (r *repository) Cancel(ctx context.Context, orderUuid uuid.UUID) error {
	query := `
        UPDATE orders
        SET order_status = $1,
            updated_at = now()
        WHERE uuid = $2
    `
	res, err := r.pool.Exec(ctx, query, model.OrderStatusCANCELLED, orderUuid)
	if res.RowsAffected() != 1 {
		if res.RowsAffected() == 0 {
			return model.OrderNotFound
		}
		log.Printf("Mismatch count affected rows for canceled")
	}
	if err != nil {
		log.Printf("failed to cancel order: %v\n", err)
		return err
	}

	return nil
}
