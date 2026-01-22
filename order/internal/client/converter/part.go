package converter // Package converter client.converter

import (
	"github.com/DeDevir/go_homework/order/internal/model"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
)

func PartFilterModelToProto(filter model.PartFilter) *inventoryV1.PartsFilter {
	return &inventoryV1.PartsFilter{
		Uuids:                 filter.UUIDS,
		Names:                 filter.Names,
		Categories:            partCategoriesModelToProto(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func partCategoriesModelToProto(categories []model.Category) []inventoryV1.Category {
	protoCategories := make([]inventoryV1.Category, len(categories))
	for _, v := range categories {
		protoCategories = append(protoCategories, partCategoryModelToProto(v))
	}
	return protoCategories
}

func partCategoryModelToProto(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryEngine:
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

func partCategoryProtoToModel(category inventoryV1.Category) model.Category {
	switch category {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func PartsProtoToModel(parts []*inventoryV1.Part) []*model.Part {
	partsModels := make([]*model.Part, 0, len(parts))
	for _, v := range parts {
		partsModels = append(partsModels, partProtoToModel(v))
	}
	return partsModels
}

func partProtoToModel(part *inventoryV1.Part) *model.Part {
	return &model.Part{
		Uuid:          uuid.MustParse(part.Uuid),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: uint64(part.StockQuantity),
		Category:      partCategoryProtoToModel(part.Category),
		Dimension:     partDimensionProtoToModel(part.Dimensions),
		Manufacturer:  partManufacturerProtoToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      partMetadatasProtoToModel(part.Metadata),
	}
}
func partDimensionProtoToModel(dimension *inventoryV1.Dimensions) *model.Dimension {
	if dimension == nil {
		return nil
	}
	return &model.Dimension{
		Width:  dimension.Width,
		Height: dimension.Height,
		Length: dimension.Length,
		Weight: dimension.Weight,
	}
}
func partManufacturerProtoToModel(manufacturer *inventoryV1.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}
func partMetadatasProtoToModel(metadata map[string]*inventoryV1.MetadataValue) map[string]*model.Metadata {
	resp := make(map[string]*model.Metadata)
	for i, v := range metadata {
		metadataPart := partMetadataProtoToModel(v)
		resp[i] = metadataPart
	}
	if len(resp) == 0 {
		return nil
	}
	return resp
}
func partMetadataProtoToModel(value *inventoryV1.MetadataValue) *model.Metadata {
	if value == nil {
		return nil
	}
	switch x := value.Value.(type) {
	case *inventoryV1.MetadataValue_DoubleValue:
		return &model.Metadata{
			Kind: model.MetadataKindDouble,
			Value: model.MetadataValue{
				Double: &x.DoubleValue,
			},
		}
	case *inventoryV1.MetadataValue_Int64Value:
		return &model.Metadata{
			Kind: model.MetadataKindInt64,
			Value: model.MetadataValue{
				Int: &x.Int64Value,
			},
		}
	case *inventoryV1.MetadataValue_BoolValue:
		return &model.Metadata{
			Kind: model.MetadataKindBool,
			Value: model.MetadataValue{
				Bool: &x.BoolValue,
			},
		}
	case *inventoryV1.MetadataValue_StringValue:
		return &model.Metadata{
			Kind: model.MetadataKindString,
			Value: model.MetadataValue{
				String: &x.StringValue,
			},
		}
	default:
		return nil
	}

}

// Еще раз понять имплементацию методов, когда мы пишем что это интерфейс и как точно понимаем что его реализовали.
// Пересмотреть структуры понять что означает * в полях структур, енамы как работают, свои типы.
// Понять как правильно работать с metadata, просмотреть на примере, что такое kind, value
