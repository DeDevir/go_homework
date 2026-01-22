package v1

import (
	"context"
	"github.com/DeDevir/go_homework/order/internal/client/converter"
	"github.com/DeDevir/go_homework/order/internal/model"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartFilter) ([]*model.Part, error) {
	resp, err := c.generatedClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: converter.PartFilterModelToProto(filter),
	})

	if err != nil {
		return nil, err
	}

	if len(resp.Parts) == 0 {
		return nil, model.SearchablePartsNotFound
	}

	return converter.PartsProtoToModel(resp.Parts), nil
}
