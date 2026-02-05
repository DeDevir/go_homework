package main

import (
	"context"
	"fmt"
	"github.com/DeDevir/go_homework/payment/internal/app"
	"github.com/DeDevir/go_homework/payment/internal/config"
	"github.com/DeDevir/go_homework/platform/pkg/closer"
	"github.com/DeDevir/go_homework/platform/pkg/logger"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
	"time"
)

const (
	configPath = "./deploy/compose/payment/.env"
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %v", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer appCancel()
	defer gracefulShutdown()

	a, err := app.New(appCtx)

	go func() {
		err = a.Run(appCtx)
		if err != nil {
			panic(fmt.Errorf("failed to run app: %v", err))
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
