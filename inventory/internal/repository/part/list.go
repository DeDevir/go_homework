package part

import (
	"context"
	"github.com/DeDevir/go_homework/inventory/internal/model"
	"github.com/DeDevir/go_homework/inventory/internal/repository/converter"
	repoModel "github.com/DeDevir/go_homework/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func (r *repository) List(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	repoFilter := converter.PartsFilterModelToInfo(filter)
	mongoFilter := buildMongoFilter(repoFilter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Failed to close mongo cursor")
		}
	}(cursor, ctx)

	var parts []*model.Part
	if err := cursor.All(ctx, &parts); err != nil {
		return nil, err
	}

	if len(parts) == 0 {
		return nil, model.PartNotFound
	}

	return parts, nil
}

func buildMongoFilter(f *repoModel.Filter) bson.D {
	filter := bson.D{}

	if len(f.UUIDs) > 0 {
		filter = append(filter, bson.E{
			Key:   "uuid",
			Value: bson.M{"$in": keys(f.UUIDs)},
		})
	}

	if len(f.Names) > 0 {
		filter = append(filter, bson.E{
			Key:   "name",
			Value: bson.M{"$in": keys(f.Names)},
		})
	}

	if len(f.Tags) > 0 {
		filter = append(filter, bson.E{
			Key:   "tags",
			Value: bson.M{"$in": keys(f.Tags)},
		})
	}

	if len(f.Categories) > 0 {
		filter = append(filter, bson.E{
			Key:   "category",
			Value: bson.M{"$in": keys(f.Categories)},
		})
	}

	if len(f.ManufacturerCountries) > 0 {
		filter = append(filter, bson.E{
			Key:   "manufacturer.country",
			Value: bson.M{"$in": keys(f.ManufacturerCountries)},
		})
	}

	return filter
}

func keys[T comparable](m map[T]struct{}) []T {
	out := make([]T, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
