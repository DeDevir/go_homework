package part

import (
	repoModel "github.com/DeDevir/go_homework/inventory/internal/repository/model"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"sync"
	"time"
)

type repository struct {
	mu   sync.RWMutex
	data map[string]*repoModel.Part
}

func NewRepository() *repository {
	uuidExample := uuid.New().String()
	mapExampleParts := map[string]*repoModel.Part{
		uuidExample: {
			Uuid:          uuidExample,
			Name:          "Wing(Bosh)",
			Description:   "Крыло от боинга по скидке",
			Price:         20000,
			StockQuantity: 2,
			Category:      repoModel.CategoryWing,
			Dimension: &repoModel.Dimension{
				Width:  40,
				Height: 30,
				Length: 100,
				Weight: 214,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Bosh",
				Country: "Deu",
				Website: "https://bosh.du",
			},
			Tags: []string{
				"wing", "bosh", "boeing",
			},
			Metadata: map[string]*repoModel.Metadata{
				"color": {
					Kind:   repoModel.MetadataKindString,
					String: lo.ToPtr("white"),
				},
			},
			CreatedAt: lo.ToPtr(time.Now()),
			UpdatedAt: lo.ToPtr(time.Now()),
		},
	}
	return &repository{
		data: mapExampleParts,
	}
}
