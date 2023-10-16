package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	internalserver "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/server"
	internalgrpc "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/server/http"
)

var (
	configFile  string
	showVersion bool

	cmd = &cobra.Command{
		Use:   "calendar",
		Short: "API for calendar app",
		RunE: func(*cobra.Command, []string) error {
			return run()
		},
	}
)

func init() {
	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show app version")
	cmd.Flags().StringVar(&configFile, "config", "/etc/calendar/config.toml", "path to config file")
}

func run() error {
	if showVersion {
		printVersion()
		return nil
	}

	ctx := context.Background()

	cfg, err := initConfig(configFile)
	if err != nil {
		return err
	}
	logger := setupLogger()

	stores, err := initStorages(ctx, cfg)
	if err != nil {
		return err
	}
	defer closeStorages(ctx, stores)

	calendar := app.New(logger, stores)

	servers := []internalserver.Server{
		internalhttp.NewServer(logger, calendar, &cfg.HTTPServer),
		internalgrpc.NewServer(logger, calendar, &cfg.GRPCServer),
	}

	notifyCtx, notifyCancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer notifyCancel()

	logger.Info("starting up calendar...")
	for _, server := range servers {
		go func(srv internalserver.Server) {
			logger.Info(fmt.Sprintf("%s server running on %s", srv.GetType(), srv.GetAddress()))
			if err := srv.Start(notifyCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("failed to start %s server: %s", srv.GetType(), err)
			}
		}(server)
	}

	<-notifyCtx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer shutdownCancel()

	logger.Info("shutting down calendar...")
	for _, server := range servers {
		if err := server.Stop(shutdownCtx); err != nil {
			logger.Error(fmt.Sprintf("failed to stop %s server: %s", server.GetType(), err))
		}
		logger.Info(fmt.Sprintf("%s server stopped", server.GetType()))
	}

	return nil
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
