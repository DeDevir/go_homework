package part

import (
	"context"
	"github.com/DeDevir/go_homework/inventory/internal/model"
	"github.com/DeDevir/go_homework/inventory/internal/repository/converter"
	repoModel "github.com/DeDevir/go_homework/inventory/internal/repository/model"
	"github.com/go-faster/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	var part repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.PartNotFound
		}
		return nil, err
	}

	return converter.PartInfoToModel(&part), nil
}
