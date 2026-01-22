package v1

import (
	"context"
	"github.com/DeDevir/go_homework/inventory/internal/converter"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, request *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.partService.List(ctx, converter.PartFiltersProtoToModel(request.Filter))

	if err != nil {
		return nil, err
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartsModelToProto(parts),
	}, nil
}
