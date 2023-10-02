package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger *slog.Logger
	stores *storage.Storages
}

func New(logger *slog.Logger, storage *storage.Storages) *App {
	return &App{
		logger: logger,
		stores: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	return a.stores.Events.Create(ctx, event)
}

func (a *App) GetEvents(ctx context.Context) ([]models.Event, error) {
	return a.stores.Events.GetMany(ctx)
}

func (a *App) GetEvent(ctx context.Context, id int64) (*models.Event, error) {
	event, err := a.stores.Events.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("%w; id=%d", ErrDocumentNotFound, id)
	}

	return event, nil
}

func (a *App) UpdateEvent(ctx context.Context, id int64, event *models.Event) (*models.Event, error) {
	event, err := a.GetEvent(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("%w; id=%d", ErrDocumentNotFound, id)
	}

	return a.stores.Events.Update(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id int64) error {
	event, err := a.GetEvent(ctx, id)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("%w; id=%d", ErrDocumentNotFound, id)
	}

	return a.stores.Events.Delete(ctx, id)
}
