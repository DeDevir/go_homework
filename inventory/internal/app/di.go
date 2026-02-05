package app

import (
	"context"
	"fmt"
	inventoryV1API "github.com/DeDevir/go_homework/inventory/internal/api/inventory/v1"
	"github.com/DeDevir/go_homework/inventory/internal/config"
	"github.com/DeDevir/go_homework/inventory/internal/repository"
	partRepository "github.com/DeDevir/go_homework/inventory/internal/repository/part"
	"github.com/DeDevir/go_homework/inventory/internal/service"
	partService "github.com/DeDevir/go_homework/inventory/internal/service/part"
	"github.com/DeDevir/go_homework/platform/pkg/closer"
	"github.com/DeDevir/go_homework/platform/pkg/logger"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	partService service.PartService

	partRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDIContainer() *diContainer { return &diContainer{} }

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = partService.NewService(d.PartRepository(ctx))
	}

	return d.partService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepository == nil {
		d.partRepository = partRepository.NewRepository(d.MongoDBHandle(ctx))
	}

	return d.partRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		logger.Info(ctx, config.AppConfig().Mongo.URI())
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}
