package payment

import (
	"context"
	"github.com/google/uuid"
)

func (*service) Pay(_ context.Context, _ string, _ string, _ string) (string, error) {
	transactionUUID := uuid.NewString()

	return transactionUUID, nil
}
