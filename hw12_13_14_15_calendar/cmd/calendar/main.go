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
		Run: func(*cobra.Command, []string) {
			run()
		},
	}
)

func init() {
	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show app version")
	cmd.Flags().StringVar(&configFile, "config", "/etc/calendar/config.toml", "path to config file")
}

func run() {
	if showVersion {
		printVersion()
		return
	}

	ctx := context.Background()

	config, err := initConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	logger := setupLogger()

	stores, err := initStorages(ctx, &config.Storage)
	if err != nil {
		log.Fatal(err)
	}
	defer closeStorages(ctx, stores)

	calendar := app.New(logger, stores)

	server := internalhttp.NewServer(logger, calendar, &config.Server)

	notifyCtx, notifyCancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer notifyCancel()

	go func() {
		logger.Info("calendar is running on " + config.Server.Address)
		if err := server.Start(notifyCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start http server: " + err.Error())
		}
	}()

	<-notifyCtx.Done()
	logger.Info("shutting down calendar...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer shutdownCancel()
	if err := server.Stop(shutdownCtx); err != nil {
		logger.Error("failed to stop http server: " + err.Error())
	}
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
