package main

import (
	"context"
	"errors"
	"fmt"
	orderV1 "github.com/DeDevir/go_homework/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/DeDevir/go_homework/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*orderV1.OrderDto),
	}
}

type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func newOrderHandler(inventoryClient inventoryV1.InventoryServiceClient, paymentClient paymentV1.PaymentServiceClient, storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func (o *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {

	uuidOrder := params.OrderUUID

	order := o.storage.GetOrderByUUID(uuidOrder)
	if order == nil {
		return &orderV1.NotFoundError{Code: 404, Message: "Order not found"}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{Code: 409, Message: "Conflict"}, nil
	}

	o.storage.CancelOrder(order)

	return &orderV1.CancelOrderNoContent{}, nil
}

func (o *OrderHandler) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	uuidOrder := params.OrderUUID
	order := o.storage.GetOrderByUUID(uuidOrder)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order not found",
		}, nil
	}
	return order, nil
}

func (o *OrderHandler) NewOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.NewOrderRes, error) {
	partsUUID := req.PartUuids
	userUUID := req.UserUUID
	var amount float64

	partsUUIDs := make([]string, 0, len(partsUUID))

	for _, u := range partsUUID {
		partsUUIDs = append(partsUUIDs, u.String())
	}

	inventoryServiceListPartsRequest := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: partsUUIDs,
		},
	}

	response, err := o.inventoryClient.ListParts(ctx, inventoryServiceListPartsRequest)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}
	if len(partsUUID) != len(response.Parts) {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Part not found",
		}, nil
	}

	for _, part := range response.Parts {
		amount += part.Price
	}

	order := o.storage.CreateOrder(userUUID, partsUUID, amount)

	return &orderV1.CreateOrderResponse{
		UUID:       orderV1.NewOptUUID(order.OrderUUID),
		TotalPrice: order.TotalPrice,
	}, nil
}

func mapPaymentMethodToGRPC(
	m orderV1.PaymentMethod,
) (paymentV1.PaymentMethod, error) {

	switch m {
	case orderV1.PaymentMethodPAYMENTMETHODCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case orderV1.PaymentMethodPAYMENTMETHODSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case orderV1.PaymentMethodPAYMENTMETHODCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	case orderV1.PaymentMethodPAYMENTMETHODUNKNOWN:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED, nil
	default:
		return 0, fmt.Errorf("unknown payment method: %s", m)
	}
}

func (o *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := o.storage.GetOrderByUUID(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order not found",
		}, nil
	}
	if order.Status != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order not found",
		}, nil
	}
	payMethod, err := mapPaymentMethodToGRPC(req.PaymentMethod)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	resp, err := o.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: payMethod,
	})
	if err != nil || resp == nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	transactionUUID := resp.TransactionUuid
	o.storage.PaidOrder(order, uuid.MustParse(transactionUUID), string(req.PaymentMethod))
	return &orderV1.PayOrderResponse{
		TransactionUUID: uuid.MustParse(transactionUUID),
	}, nil
}

func (o *OrderHandler) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func (s *OrderStorage) GetOrderByUUID(uuid uuid.UUID) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderStorage) PaidOrder(order *orderV1.OrderDto, transactionUUID uuid.UUID, paymentMethod string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order.TransactionUUID = orderV1.NewOptUUID(transactionUUID)
	order.Status = orderV1.OrderStatusPAID
	order.PaymentMethod = orderV1.NewOptString(paymentMethod)

}

func (s *OrderStorage) CreateOrder(userUUID uuid.UUID, partsUUID []uuid.UUID, amount float64) *orderV1.OrderDto {
	order := &orderV1.OrderDto{
		OrderUUID:       uuid.New(),
		UserUUID:        userUUID,
		PartUuids:       partsUUID,
		TotalPrice:      amount,
		TransactionUUID: orderV1.OptUUID{},
		PaymentMethod:   orderV1.OptString{},
		Status:          orderV1.OrderStatusPENDINGPAYMENT,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.OrderUUID] = order
	return order
}

func (s *OrderStorage) CancelOrder(order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order.Status = orderV1.OrderStatusCANCELLED
}

func main() {

	inventoryConn, err := grpc.NewClient(
		net.JoinHostPort("localhost", inventoryServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Printf("Err start conn orderGRPC %v", err)
		return
	}

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	payConn, err := grpc.NewClient(
		net.JoinHostPort("localhost", paymentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Printf("Err start conn orderGRPC %v", err)
		return
	}

	paymentClient := paymentV1.NewPaymentServiceClient(payConn)

	handler := newOrderHandler(inventoryClient, paymentClient, NewOrderStorage())
	orderServer, err := orderV1.NewServer(handler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
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
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

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
