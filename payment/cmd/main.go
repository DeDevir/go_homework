package main

import (
	"context"
	"fmt"
	paymentV1 "github.com/DeDevir/go_homework/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	port = 50052
)

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (p *paymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	orderUUID := req.OrderUuid
	if _, err := uuid.Parse(orderUUID); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "orderUUID not corrected %v", err)
	}

	userUUID := req.UserUuid
	if _, err := uuid.Parse(userUUID); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "userUUID not corrected %v", err)
	}

	paymentMethod := req.PaymentMethod
	if paymentMethod == paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "payment method cannot be unknown")
	}

	transactionUUID := uuid.New().String()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s\n", transactionUUID)
	return &paymentV1.PayOrderResponse{TransactionUuid: transactionUUID}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Failed to listen %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("Failed to close listener %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	service := &paymentService{}

	paymentV1.RegisterPaymentServiceServer(s, service)
	reflection.Register(s)

	go func() {
		log.Printf(" grpc server start listen on %v port", port)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("Failed to serve %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
