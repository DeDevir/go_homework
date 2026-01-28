package part

import (
	"context"
	repoModel "github.com/DeDevir/go_homework/inventory/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	collection := initCollectionPart(db)

	return &repository{
		collection: collection,
	}
}

func initCollectionPart(db *mongo.Database) *mongo.Collection {
	collection := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "tags", Value: 1}}, // multikey автоматически
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "manufacturer.country", Value: 1}},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}
	//for _, idxName := range indexNames {
	//	log.Printf("Создан индекс: %s\n", idxName)
	//}
	insertExamplePart(collection)
	return collection
}

func insertExamplePart(collection *mongo.Collection) {
	uuidExample := uuid.New().String()
	examplePart := repoModel.Part{
		Uuid:          uuidExample,
		Name:          gofakeit.AppName(),
		Description:   gofakeit.AirlineFlightNumber(),
		Price:         float64(gofakeit.IntRange(100, 50000)),
		StockQuantity: uint64(gofakeit.IntRange(3, 200)),
		Category:      repoModel.Category(gofakeit.IntRange(0, 4)),
		Dimension: &repoModel.Dimension{
			Width:  40,
			Height: 30,
			Length: 100,
			Weight: 214,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
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
	}

	_, err := collection.InsertOne(context.Background(), &examplePart)
	if err != nil {
		log.Printf("Ошибка при вставке детали %v\n", err)
	}
}
