package env

import (
	"github.com/caarlos0/env/v11"
	"net"
	"time"
)

type orderHttpEnvConfig struct {
	Host        string `env:"HTTP_HOST,required"`
	Port        string `env:"HTTP_PORT,required"`
	ReadTimeout string `env:"HTTP_READ_TIMEOUT,required"`
}

type orderHttpConfig struct {
	raw orderHttpEnvConfig
}

func NewOrderHTTPConfig() (*orderHttpConfig, error) {
	var raw orderHttpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderHttpConfig{raw: raw}, nil
}

func (cfg *orderHttpConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *orderHttpConfig) ReadTimeout() time.Duration {
	readTimeout, err := time.ParseDuration(cfg.raw.ReadTimeout)
	if err != nil {
		return 5 * time.Second
	}
	return readTimeout
}
