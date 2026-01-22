package main

import (
	"fmt"
	inventoryApi "github.com/DeDevir/go_homework/inventory/internal/api/inventory/v1"
	partRepository "github.com/DeDevir/go_homework/inventory/internal/repository/part"
	partService "github.com/DeDevir/go_homework/inventory/internal/service/part"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
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
	listener, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", port),
	)
	if err != nil {
		log.Printf("Error launch listener: %v", err)
		return
	}

	partRepositoryLay := partRepository.NewRepository()
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
