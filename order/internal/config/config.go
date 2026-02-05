package config

import (
	"github.com/DeDevir/go_homework/order/internal/config/env"
	"github.com/joho/godotenv"
	"os"
)

var appConfig *config

type config struct {
	Logger              LoggerConfig
	InventoryGRPCClient InventoryClientGRPCConfig
	PaymentGRPCClient   PaymentClientGRPCConfig
	OrderHTTP           OrderHTTPConfig
	Postgres            PostgresConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	invGRPCClientCfg, err := env.NewInventoryGRPCClientConfig()
	if err != nil {
		return err
	}

	paymentGRPCClientCfg, err := env.NewPaymentGRPCClientConfig()
	if err != nil {
		return err
	}

	orderHttpCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:              loggerCfg,
		InventoryGRPCClient: invGRPCClientCfg,
		PaymentGRPCClient:   paymentGRPCClientCfg,
		OrderHTTP:           orderHttpCfg,
		Postgres:            postgresCfg,
	}
	return nil
}

func AppConfig() *config { return appConfig }
