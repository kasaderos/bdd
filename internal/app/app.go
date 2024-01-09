package app

import (
	"bdd/config"
	"bdd/internal/handler"
	"bdd/internal/repository"
	"bdd/internal/service"
	"bdd/pkg/logger"
	"bdd/pkg/postgres"
	"bdd/pkg/template"
	"context"
	"log"
	"log/slog"
	"net/http"
	"sync"
)

type App struct {
	ctx context.Context

	service    *service.Service
	repo       *repository.Postgres
	httpServer *http.Server
	conf       *config.Config

	wg        sync.WaitGroup
	cancelCtx context.CancelFunc
	errC      chan error
}

func New() *App {
	// Global context
	ctx, cancel := context.WithCancel(context.Background())

	// Config
	conf := config.Load()

	// Logger
	slog.SetDefault(logger.New(conf.LogLevel))

	// Database
	dbpool, err := postgres.Open(ctx, conf.DatabaseConfig.PostgresDSN())
	if err != nil {
		log.Fatal(err)
	}

	// Repository
	repo := repository.NewPostgresRepo(dbpool)

	// Templates
	templManager := template.NewManager()
	if err := templManager.ParseTemplates(
		conf.ServiceConfig.Templates,
	); err != nil {
		log.Fatal(err)
	}

	// Service
	srv := service.New(
		repo,
		conf.ServiceConfig,
	)

	// Handlers
	httpServer := new(http.Server)
	httpServer.Addr = ":" + conf.HTTPListenPort
	httpServer.Handler = handler.NewHTTPHandler(
		srv,
		templManager,
	)

	return &App{
		ctx:        ctx,
		service:    srv,
		repo:       repo,
		httpServer: httpServer,
		conf:       conf,
		cancelCtx:  cancel,
		errC:       make(chan error, 1),
	}
}

func (a *App) Run() {
	a.startWorker(a.initShutdown)
	a.startWorker(a.startServer)

	a.wait()
}

func (a *App) startWorker(f func(context.Context)) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		f(a.ctx)
	}()
}

func (a *App) wait() {
	go func() {
		a.wg.Wait()
		close(a.errC)
	}()

	for err := range a.errC {
		if err != nil {
			slog.Error("app", "msg", err.Error())
		}
	}
}
