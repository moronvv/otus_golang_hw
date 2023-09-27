package internalhttp

import (
	"context"
	"net/http"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
)

type Logger interface { // TODO
}

type Application interface { // TODO
}

type Server struct { // TODO
	ctx    context.Context
	server *http.Server
	logger Logger
	app    Application
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloRoute)

	return mux
}

func NewServer(logger Logger, app Application, cfg *config.ServerConf) *Server {
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: setupRoutes(),
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
