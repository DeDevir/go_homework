package app

import (
	"context"
	"errors"
	"github.com/DeDevir/go_homework/order/internal/config"
	"github.com/DeDevir/go_homework/order/internal/migrator"
	"github.com/DeDevir/go_homework/platform/pkg/closer"
	"github.com/DeDevir/go_homework/platform/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type App struct {
	diContainer *diContainer
	router      *chi.Mux
	httpServer  *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHttpServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initMigrator,
		a.initRouter,
		a.initServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initDi(ctx context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	return logger.Init(config.AppConfig().Logger.Level(), config.AppConfig().Logger.AsJson())
}

func (a *App) initMigrator(ctx context.Context) error {
	poolCfg, err := pgxpool.ParseConfig(config.AppConfig().Postgres.URI())
	if err != nil {
		return err
	}

	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*poolCfg.ConnConfig), config.AppConfig().Postgres.MigrationDirectory())
	err = migratorRunner.Up(ctx)
	if err != nil {
		logger.Error(ctx, "migrator up failed", zap.Error(err))
		return err
	}
	logger.Info(ctx, "migrator up done")
	return nil
}

func (a *App) initRouter(ctx context.Context) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", a.diContainer.OrderServer(ctx))
	a.router = r
	return nil
}

func (a *App) initServer(_ context.Context) error {
	server := &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           a.router,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(), // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}
	a.httpServer = server
	return nil
}

func (a *App) runHttpServer(ctx context.Context) error {
	logger.Info(ctx, "starting http server ", zap.String("address", config.AppConfig().OrderHTTP.Address()))
	closer.AddNamed("Http Server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})
	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
