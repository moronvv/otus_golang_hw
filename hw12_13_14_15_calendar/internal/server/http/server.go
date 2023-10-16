package internalhttp

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	internalserver "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/server"
	internalhttproutes "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/server/http/routes"
)

type server struct {
	srv    *http.Server
	logger *slog.Logger
	cfg    *config.HTTPServerConf
	app    app.App
}

func NewServer(logger *slog.Logger, app app.App, cfg *config.HTTPServerConf) internalserver.Server {
	router := internalhttproutes.SetupRoutes(app)
	router.Use(newLoggerMiddleware(logger).Middleware)

	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.RequestTimeout,
	}

	return &server{
		srv:    srv,
		logger: logger,
		cfg:    cfg,
		app:    app,
	}
}

func (s *server) GetType() string {
	return "HTTP"
}

func (s *server) GetAddress() string {
	return s.cfg.Address
}

func (s *server) Start(context.Context) error {
	return s.srv.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
