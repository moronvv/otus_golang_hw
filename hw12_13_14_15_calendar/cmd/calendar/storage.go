package main

import (
	"context"
	"fmt"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/memory"
)

func initStorages(ctx context.Context, storageConf *config.StorageConf) (*storage.Storages, error) {
	var store storage.Storage
	var eventStorage storage.EventStorage

	switch storageConf.Type {
	case "in-memory":
		store = memorystorage.NewStorage()
		if err := store.Connect(ctx); err != nil {
			return nil, err
		}
		eventStorage = memorystorage.NewEventStorage(store.(*memorystorage.InMemoryStorage))
	case "sql":
		// TODO: implement sql
		return nil, fmt.Errorf("not implemented")
	default:
		return nil, fmt.Errorf("unsupported storage type %s", storageConf.Type)
	}

	return &storage.Storages{
		DB:     store,
		Events: eventStorage,
	}, nil
}

func closeStorages(ctx context.Context, storages *storage.Storages) error {
	return storages.DB.Close(ctx)
}
