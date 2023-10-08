package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
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

	server := internalhttp.NewServer(logger, calendar, cfg)

	notifyCtx, notifyCancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer notifyCancel()

	go func() {
		logger.Info("calendar is running on " + cfg.Server.Address)
		if err := server.Start(notifyCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start http server: " + err.Error())
		}
	}()

	<-notifyCtx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer shutdownCancel()

	logger.Info("shutting down calendar...")
	if err := server.Stop(shutdownCtx); err != nil {
		logger.Error("failed to stop http server: " + err.Error())
		return err
	}

	return nil
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
