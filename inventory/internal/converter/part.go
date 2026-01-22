package converter

import (
	"github.com/DeDevir/go_homework/inventory/internal/model"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PartFiltersProtoToModel(Proto *inventoryV1.PartsFilter) *model.PartsFilter {
	return &model.PartsFilter{
		Uuid:                  Proto.Uuids,
		Names:                 Proto.Names,
		Categories:            categoriesProtoToModels(Proto.Categories),
		ManufacturerCountries: Proto.ManufacturerCountries,
		Tags:                  Proto.Tags,
	}
}

func categoriesProtoToModels(categories []inventoryV1.Category) []model.Category {
	categoriesModel := make([]model.Category, len(categories))
	for i, ctg := range categories {
		categoriesModel[i] = categoryProtoToModel(ctg)
	}
	return categoriesModel
}

func categoryProtoToModel(category inventoryV1.Category) model.Category {
	switch category {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryENGINE
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUNKNOWN
	}
}

func categoryModelToProto(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryENGINE:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNKNOWN_UNSPECIFIED
	}
}

func PartsModelToProto(parts []*model.Part) []*inventoryV1.Part {
	partsProto := make([]*inventoryV1.Part, len(parts))

	for i, part := range parts {
		partsProto[i] = PartModelToProto(part)
	}

	return partsProto
}

func PartModelToProto(part *model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: int64(part.StockQuantity),
		Category:      categoryModelToProto(part.Category),
		Dimensions:    dimensionModelToProto(part.Dimension),
		Manufacturer:  manufacturerModelToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataModelToProto(part.Metadata),
		CreatedAt:     timestamppb.New(*part.CreatedAt),
		UpdatedAt:     timestamppb.New(*part.UpdatedAt),
	}
}

func dimensionModelToProto(dimension *model.Dimension) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: dimension.Length,
		Width:  dimension.Width,
		Height: dimension.Height,
		Weight: dimension.Weight,
	}
}

func manufacturerModelToProto(manufacturer *model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func metadataModelToProto(m map[string]*model.MetadataValue) map[string]*inventoryV1.MetadataValue {
	if m == nil {
		return nil
	}
	res := make(map[string]*inventoryV1.MetadataValue, len(m))
	for i, v := range m {
		switch v.Kind {
		case model.MetadataKindBool:
			if v.Bool == nil {
				continue
			}
			res[i] = &inventoryV1.MetadataValue{Value: &inventoryV1.MetadataValue_BoolValue{BoolValue: *v.Bool}}
		case model.MetadataKindDouble:
			if v.Double == nil {
				continue
			}
			res[i] = &inventoryV1.MetadataValue{Value: &inventoryV1.MetadataValue_DoubleValue{DoubleValue: *v.Double}}
		case model.MetadataKindInt64:
			if v.Int64 == nil {
				continue
			}
			res[i] = &inventoryV1.MetadataValue{Value: &inventoryV1.MetadataValue_Int64Value{Int64Value: *v.Int64}}
		case model.MetadataKindString:
			if v.String == nil {
				continue
			}
			res[i] = &inventoryV1.MetadataValue{Value: &inventoryV1.MetadataValue_StringValue{StringValue: *v.String}}
		default:
			continue
		}
	}
	return res
}
