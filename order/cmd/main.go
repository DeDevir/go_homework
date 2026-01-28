package main

import (
	"context"
	"errors"
	orderApi "github.com/DeDevir/go_homework/order/internal/api/order/v1"
	inventoryClientGrpc "github.com/DeDevir/go_homework/order/internal/client/grpc/inventory/v1"
	paymentClientGrpc "github.com/DeDevir/go_homework/order/internal/client/grpc/payment/v1"
	"github.com/DeDevir/go_homework/order/internal/migrator"
	orderRepository "github.com/DeDevir/go_homework/order/internal/repository/order"
	orderService "github.com/DeDevir/go_homework/order/internal/service/order"
	orderV1 "github.com/DeDevir/go_homework/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/DeDevir/go_homework/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	inventoryServicePort = "50051"
	paymentServicePort   = "50052"
)

func main() {
	ctx := context.Background()

	invetoryConn, err := grpc.NewClient(
		net.JoinHostPort("localhost", inventoryServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Error start connection inventory grpc")
		return
	}
	inventoryServiceClientGrpc := inventoryV1.NewInventoryServiceClient(invetoryConn)
	inventoryClient := inventoryClientGrpc.NewClient(inventoryServiceClientGrpc)

	paymentConn, err := grpc.NewClient(
		net.JoinHostPort("localhost", paymentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Error start connection payment grpc")
		return
	}

	dbURI := "postgres://order-service-user:order-service-password@localhost:5432/order-service"
	con, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		con.Close()
	}()

	// Проверяем, что соединение с базой установлено
	err = con.Ping(ctx)
	if err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		return
	}

	migrationsDir := "order/migrations"
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*con.Config().ConnConfig.Copy()), migrationsDir)
	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Ошибка миграции базы данных: %v\n", err)
		return
	}

	paymentServiceClientGrpc := paymentV1.NewPaymentServiceClient(paymentConn)
	paymentClient := paymentClientGrpc.NewClient(paymentServiceClientGrpc)
	repository := orderRepository.NewRepository(con)
	service := orderService.NewService(repository, inventoryClient, paymentClient)
	api := orderApi.NewApi(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("Error start order server")
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("Запускаем http сервер order на порту: %v\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Ошибка запуска order сервера")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
