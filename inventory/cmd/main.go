package main

import (
	"context"
	"fmt"
	inventoryApi "github.com/DeDevir/go_homework/inventory/internal/api/inventory/v1"
	partRepository "github.com/DeDevir/go_homework/inventory/internal/repository/part"
	partService "github.com/DeDevir/go_homework/inventory/internal/service/part"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	port = 50051
)

func main() {
	ctx := context.Background()
	listener, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", port),
	)
	if err != nil {
		log.Printf("Error launch listener: %v", err)
		return
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://inventory-service-user:inventory-service-password@localhost:27017"))
	if err != nil {
		log.Printf("Ошибка подключения к mongo database: %v\n", err)
	}
	mongoDb := client.Database("inventory-service") // docker

	partRepositoryLay := partRepository.NewRepository(mongoDb)
	partServiceLay := partService.NewService(partRepositoryLay)
	api := inventoryApi.NewAPI(partServiceLay)

	server := grpc.NewServer()

	inventoryV1.RegisterInventoryServiceServer(server, api)

	reflection.Register(server)

	go func() {
		log.Printf("GRPC Inventory Server start to listen on %v port", port)
		err = server.Serve(listener)
		if err != nil {
			log.Printf("Failed to serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	server.GracefulStop()
	log.Println("✅ Server stopped")
}
