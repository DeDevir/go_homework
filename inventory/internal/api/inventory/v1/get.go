package v1

import (
	"context"
	"github.com/DeDevir/go_homework/inventory/internal/converter"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
	"log"
)

func (a *api) GetPart(ctx context.Context, request *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.partService.Get(ctx, request.Uuid)
	if err != nil {
		log.Printf("failed to get part - %v", err)
		return nil, err
	}

	return &inventoryV1.GetPartResponse{Part: converter.PartModelToProto(part)}, nil
}
