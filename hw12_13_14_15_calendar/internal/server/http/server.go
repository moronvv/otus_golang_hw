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
	srv *http.Server
	cfg *config.HTTPServerConf
	app app.App
}

func NewBaseServer(logger *slog.Logger, app app.App, cfg *config.HTTPServerConf) *http.Server {
	router := internalhttproutes.SetupRoutes(app)
	router.Use(newLoggerMiddleware(logger).Middleware)
	router.Use(AuthMiddleware)

	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.RequestTimeout,
	}

	return srv
}

func NewServer(
	logger *slog.Logger,
	app app.App,
	cfg *config.HTTPServerConf,
	baseSrv *http.Server,
) internalserver.Server {
	if baseSrv == nil {
		baseSrv = NewBaseServer(logger, app, cfg)
	}

	return &server{
		srv: baseSrv,
		cfg: cfg,
		app: app,
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
