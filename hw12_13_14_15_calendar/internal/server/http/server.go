package internalhttp

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
)

type Application interface { // TODO
}

type Server struct { // TODO
	ctx    context.Context
	server *http.Server
	logger *slog.Logger
	app    Application
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)

	return mux
}

func NewServer(logger *slog.Logger, app Application, cfg *config.ServerConf) *Server {
	handler := newLoggerMiddleware(logger, setupRoutes())

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: handler,
	}

	return &Server{
		server: server,
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}
