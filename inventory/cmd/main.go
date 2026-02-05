package main

import (
	"context"
	"fmt"
	"github.com/DeDevir/go_homework/inventory/internal/app"
	"github.com/DeDevir/go_homework/inventory/internal/config"
	"github.com/DeDevir/go_homework/platform/pkg/closer"
	"github.com/DeDevir/go_homework/platform/pkg/logger"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
	"time"
)

const (
	configPath = "./deploy/compose/inventory/.env"
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %v", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
