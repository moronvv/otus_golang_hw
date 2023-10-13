package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type App interface {
	CreateEvent(context.Context, *models.Event) (*models.Event, error)
	GetEvents(context.Context) ([]models.Event, error)
	GetEvent(context.Context, int64) (*models.Event, error)
	UpdateEvent(context.Context, int64, *models.Event) (*models.Event, error)
	DeleteEvent(context.Context, int64) error
}

type app struct {
	logger *slog.Logger
	stores *storage.Storages
}

func New(logger *slog.Logger, storage *storage.Storages) App {
	return &app{
		logger: logger,
		stores: storage,
	}
}

func (a *app) CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	return a.stores.Events.Create(ctx, event)
}

func (a *app) GetEvents(ctx context.Context) ([]models.Event, error) {
	return a.stores.Events.GetMany(ctx)
}

func (a *app) GetEvent(ctx context.Context, id int64) (*models.Event, error) {
	event, err := a.stores.Events.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("%w; id=%d", ErrDocumentNotFound, id)
	}

	return event, nil
}

func (a *app) UpdateEvent(ctx context.Context, id int64, event *models.Event) (*models.Event, error) {
	updated, err := a.stores.Events.Update(ctx, id, event)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, fmt.Errorf("%w; id=%d", ErrDocumentNotFound, id)
	}

	return updated, err
}

func (a *app) DeleteEvent(ctx context.Context, id int64) error {
	ok, err := a.stores.Events.Delete(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%w; id=%d", ErrDocumentNotFound, id)
	}

	return nil
}
