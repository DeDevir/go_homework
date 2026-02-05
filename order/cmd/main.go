package main

import (
	"context"
	"fmt"
	"github.com/DeDevir/go_homework/order/internal/app"
	"github.com/DeDevir/go_homework/order/internal/config"
	"github.com/DeDevir/go_homework/platform/pkg/closer"
	"github.com/DeDevir/go_homework/platform/pkg/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	configPath = "./deploy/compose/order/.env"
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %v", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()

	a, err := app.New(appCtx)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		err = a.Run(appCtx)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	<-appCtx.Done()
	gracefulShutdown()

}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
