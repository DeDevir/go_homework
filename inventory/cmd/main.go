package main

import (
	"context"
	"fmt"
	inventoryV1 "github.com/DeDevir/go_homework/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	port = 50051
)

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

type MapFiltering struct {
	UUIDS                 map[string]struct{}
	Names                 map[string]struct{}
	Categories            map[inventoryV1.Category]struct{}
	ManufacturerCountries map[string]struct{}
	Tags                  map[string]struct{}
}

func (s *InventoryService) GetPart(_ context.Context, request *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	part := s.parts[request.Uuid]
	if part == nil {
		return nil, status.Errorf(codes.NotFound, "Part not found")
	}
	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *InventoryService) ListParts(_ context.Context, request *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := request.Filter
	indexFilter := FilledMapIndexFiltering(filter)

	filteredParts := make([]*inventoryV1.Part, 0, len(s.parts))
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, part := range s.parts {
		if indexFilter.UUIDS != nil {
			if _, ok := indexFilter.UUIDS[part.Uuid]; !ok {
				continue
			}
		}

		if indexFilter.Names != nil {
			if _, ok := indexFilter.Names[part.Name]; !ok {
				continue
			}
		}

		if indexFilter.Categories != nil {
			if _, ok := indexFilter.Categories[part.Category]; !ok {
				continue
			}
		}

		if indexFilter.ManufacturerCountries != nil {
			if _, ok := indexFilter.ManufacturerCountries[part.Manufacturer.Country]; !ok {
				continue
			}
		}

		if indexFilter.Tags != nil {
			issetTag := false
			for _, t := range part.Tags {
				if _, ok := indexFilter.Tags[t]; ok {
					issetTag = true
					break
				}
			}
			if !issetTag {
				continue
			}
		}

		filteredParts = append(filteredParts, part)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: filteredParts,
	}, nil
}

func FilledMapIndexFiltering(filter *inventoryV1.PartsFilter) *MapFiltering {
	mapIndex := &MapFiltering{}
	if len(filter.Uuids) > 0 {
		mapIndex.UUIDS = make(map[string]struct{}, len(filter.Uuids))
		for _, u := range filter.Uuids {
			mapIndex.UUIDS[u] = struct{}{}
		}
	}

	if len(filter.Names) > 0 {
		mapIndex.Names = make(map[string]struct{}, len(filter.Names))
		for _, u := range filter.Names {
			mapIndex.Names[u] = struct{}{}
		}
	}

	if len(filter.Categories) > 0 {
		mapIndex.Categories = make(map[inventoryV1.Category]struct{}, len(filter.Categories))
		for _, u := range filter.Categories {
			mapIndex.Categories[u] = struct{}{}
		}
	}

	if len(filter.ManufacturerCountries) > 0 {
		mapIndex.ManufacturerCountries = make(map[string]struct{}, len(filter.ManufacturerCountries))
		for _, u := range filter.ManufacturerCountries {
			mapIndex.ManufacturerCountries[u] = struct{}{}
		}
	}

	if len(filter.Tags) > 0 {
		mapIndex.Tags = make(map[string]struct{}, len(filter.Tags))
		for _, u := range filter.Tags {
			mapIndex.Tags[u] = struct{}{}
		}
	}

	return mapIndex
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Printf("Fail to listen %v", err)
		return
	}
	defer func() {
		log.Printf("Stoping listener on %v port", port)
		err = listener.Close()
		if err != nil {
			log.Printf("Failed to close listener on %v port", port)
		}
	}()
	uuidExample := uuid.New().String()
	dimensionExample := &inventoryV1.Dimensions{
		Length: 20,
		Width:  10,
		Height: 5,
		Weight: 100,
	}
	manufacturerExample := &inventoryV1.Manufacturer{
		Name:    "Bosh",
		Country: "DU",
		Website: "https://bosh.com",
	}
	tagsExample := []string{
		"Airplane", "Deutsch", "Bosh",
	}
	metaDataExample := map[string]*inventoryV1.Value{
		"power": {Value: &inventoryV1.Value_DoubleValue{DoubleValue: 123.0}},
	}
	service := &InventoryService{
		parts: map[string]*inventoryV1.Part{
			uuidExample: {
				Uuid:          uuidExample,
				Name:          "Wing",
				Description:   "Wing from airplane",
				Price:         120,
				StockQuantity: 5,
				Category:      inventoryV1.Category_CATEGORY_WING,
				Dimensions:    dimensionExample,
				Manufacturer:  manufacturerExample,
				Tags:          tagsExample,
				Metadata:      metaDataExample,
				CreatedAt:     timestamppb.Now(),
				UpdatedAt:     timestamppb.Now(),
			},
		},
	}

	server := grpc.NewServer()

	inventoryV1.RegisterInventoryServiceServer(server, service)

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
