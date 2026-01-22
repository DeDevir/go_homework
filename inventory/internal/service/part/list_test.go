package part

import (
	"github.com/DeDevir/go_homework/inventory/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuit) TestListPartSuccess() {
	var (
		exampleUuid       = gofakeit.UUID()
		err         error = nil
		examplePart       = []*model.Part{
			{
				Uuid:          exampleUuid,
				Name:          gofakeit.PetName(),
				Description:   gofakeit.AdverbPlace(),
				Price:         gofakeit.Price(1, 200),
				StockQuantity: gofakeit.Uint64(),
				Category:      model.Category(gofakeit.UintRange(1, 7)),
				Dimension:     nil,
				Manufacturer:  nil,
				Tags:          nil,
				Metadata:      nil,
				CreatedAt:     nil,
				UpdatedAt:     nil,
			},
		}

		filter = &model.PartsFilter{
			Uuid: []string{
				exampleUuid,
			},
			Names:                 nil,
			Categories:            nil,
			ManufacturerCountries: nil,
			Tags:                  nil,
		}
	)

	s.partRepository.On("List", s.ctx, filter).Return(examplePart, err)

	resPart, err := s.service.List(s.ctx, filter)

	s.Require().NoError(err)
	s.Require().Equal(examplePart, resPart)
}
