package converter

import (
	"github.com/DeDevir/go_homework/inventory/internal/model"
	repoModel "github.com/DeDevir/go_homework/inventory/internal/repository/model"
)

func PartInfoToModel(part *repoModel.Part) *model.Part {
	return &model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      categoryInfoToModel(part.Category),
		Dimension:     dimensionInfoToModel(part.Dimension),
		Manufacturer:  manufacturerInfoToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataInfoToModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func categoryInfoToModel(category repoModel.Category) model.Category {
	return model.Category(category)
}
func dimensionInfoToModel(dimension *repoModel.Dimension) *model.Dimension {
	return &model.Dimension{
		Width:  dimension.Width,
		Height: dimension.Height,
		Length: dimension.Length,
		Weight: dimension.Weight,
	}
}
func manufacturerInfoToModel(manufacturer *repoModel.Manufacturer) *model.Manufacturer {
	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}
func metadataInfoToModel(metadata map[string]*repoModel.Metadata) map[string]*model.MetadataValue {
	res := make(map[string]*model.MetadataValue)
	for i, v := range metadata {
		res[i] = &model.MetadataValue{
			Kind:   model.MetadataKind(v.Kind),
			String: v.String,
			Int64:  v.Int64,
			Double: v.Double,
			Bool:   v.Bool,
		}
	}

	return res
}

func PartsFilterModelToInfo(filter *model.PartsFilter) *repoModel.Filter {
	uuids := make(map[string]struct{})
	tags := make(map[string]struct{})
	manufacturerCountries := make(map[string]struct{})
	names := make(map[string]struct{})
	categories := make(map[repoModel.Category]struct{})

	for _, v := range filter.Uuid {
		uuids[v] = struct{}{}
	}

	for _, v := range filter.Tags {
		tags[v] = struct{}{}
	}

	for _, v := range filter.ManufacturerCountries {
		manufacturerCountries[v] = struct{}{}
	}
	for _, v := range filter.Names {
		names[v] = struct{}{}
	}
	for _, v := range filter.Categories {
		switch v {
		case model.CategoryWing:
			categories[repoModel.CategoryWing] = struct{}{}
		case model.CategoryPorthole:
			categories[repoModel.CategoryPorthole] = struct{}{}
		case model.CategoryFuel:
			categories[repoModel.CategoryFuel] = struct{}{}
		case model.CategoryENGINE:
			categories[repoModel.CategoryEngine] = struct{}{}
		default:
			categories[repoModel.CategoryUndefined] = struct{}{}
		}
	}

	return &repoModel.Filter{
		UUIDs:                 uuids,
		Names:                 names,
		Tags:                  tags,
		Categories:            categories,
		ManufacturerCountries: manufacturerCountries,
	}
}
