package main

import (
	"context"
	"fmt"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

func initStorages(ctx context.Context, cfg *config.Config) (*storage.Storages, error) {
	var newStorageFn func(*config.DatabaseConf) storage.Storage
	var newEventStorageFn func(storage.Storage) storage.EventStorage

	switch cfg.Storage.Type {
	case "in-memory":
		newStorageFn = memorystorage.NewStorage
		newEventStorageFn = memorystorage.NewEventStorage
	case "sql":
		newStorageFn = sqlstorage.NewStorage
		newEventStorageFn = sqlstorage.NewEventStorage
	default:
		return nil, fmt.Errorf("unsupported storage type %s", cfg.Storage.Type)
	}

	store := newStorageFn(&cfg.Database)
	if err := store.Connect(ctx); err != nil {
		return nil, err
	}
	eventStorage := newEventStorageFn(store)

	return &storage.Storages{
		DB:     store,
		Events: eventStorage,
	}, nil
}

func closeStorages(ctx context.Context, storages *storage.Storages) error {
	return storages.DB.Close(ctx)
}
