package internalhttp

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	internalhttproutes "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/server/http/routes"
)

type Server struct {
	server *http.Server
	logger *slog.Logger
	app    app.App
}

func NewServer(logger *slog.Logger, app app.App, cfg *config.Config) *Server {
	router := internalhttproutes.SetupRoutes(app)
	router.Use(newLoggerMiddleware(logger).Middleware)

	server := &http.Server{
		Addr:              cfg.Server.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.Server.RequestTimeout,
	}

	return &Server{
		server: server,
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(context.Context) error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
