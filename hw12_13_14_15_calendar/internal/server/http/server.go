package internalhttp

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
)

type Application interface { // TODO
}

type Server struct { // TODO
	server *http.Server
	logger *slog.Logger
	app    Application
}

func setupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", pingHandler).Methods("GET")

	return router
}

func NewServer(logger *slog.Logger, app Application, cfg *config.Config) *Server {
	handler := newLoggerMiddleware(logger, setupRoutes())

	server := &http.Server{
		Addr:              cfg.Server.Address,
		Handler:           handler,
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
