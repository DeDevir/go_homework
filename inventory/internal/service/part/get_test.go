package part

import (
	"github.com/DeDevir/go_homework/inventory/internal/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"time"
)

func (s *ServiceSuit) TestGetSuccess() {
	var (
		searchedUuid = gofakeit.UUID()
		part         = &model.Part{
			Uuid:          searchedUuid,
			Name:          gofakeit.Unit(),
			Description:   gofakeit.Language(),
			Price:         gofakeit.Price(1, 500),
			StockQuantity: gofakeit.Uint64(),
			Category:      model.Category(gofakeit.IntRange(1, 4)),
			Dimension: &model.Dimension{
				Length: gofakeit.Float64Range(1, 50),
				Width:  gofakeit.Float64Range(1, 50),
				Height: gofakeit.Float64Range(1, 50),
				Weight: gofakeit.Float64Range(1, 50),
			},
			Manufacturer: &model.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: gofakeit.ProductAudience(),
			Metadata: map[string]*model.MetadataValue{
				gofakeit.UrlSlug(2): &model.MetadataValue{
					Kind:   model.MetadataKind("string"),
					String: lo.ToPtr("123"),
				},
			},
			CreatedAt: lo.ToPtr(time.Now()),
			UpdatedAt: lo.ToPtr(time.Now()),
		}
	)
	s.partRepository.On("Get", s.ctx, searchedUuid).Return(part, nil)

	findedPart, err := s.service.Get(s.ctx, searchedUuid)
	s.Require().NoError(err)
	s.Require().Equal(part, findedPart)
}

func (s *ServiceSuit) TestGetError() {
	var (
		repoError    = gofakeit.Error()
		searchedUuid = gofakeit.UUID()
	)

	s.partRepository.On("Get", s.ctx, searchedUuid).Return(nil, repoError)

	part, err := s.service.Get(s.ctx, searchedUuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoError)
	s.Require().Nil(part)
}

// Еще list_test сделать в сервис слое, и переходи в репозиторный слой, затем уже в сервис ордер
