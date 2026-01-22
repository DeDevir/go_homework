package v1

import (
	"github.com/DeDevir/go_homework/inventory/internal/service"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	partService service.PartService
}

func NewAPI(service service.PartService) *api {
	return &api{
		partService: service,
	}
}
