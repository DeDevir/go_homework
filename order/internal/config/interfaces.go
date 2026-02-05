package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryClientGRPCConfig interface {
	Address() string
}

type PaymentClientGRPCConfig interface {
	Address() string
}

type OrderHTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type PostgresConfig interface {
	URI() string
	MigrationDirectory() string
}
