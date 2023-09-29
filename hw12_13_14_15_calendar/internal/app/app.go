package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  *slog.Logger
	storage storage.EventStorage
}

func New(logger *slog.Logger, storage storage.EventStorage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	return a.storage.Create(ctx, event)
}

func (a *App) GetEvents(ctx context.Context) ([]models.Event, error) {
	return a.storage.GetMany(ctx)
}

func (a *App) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	event, err := a.storage.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("%w; id=%s", ErrDocumentNotFound, id)
	}

	return event, nil
}

func (a *App) UpdateEvent(ctx context.Context, id uuid.UUID, event *models.Event) (*models.Event, error) {
	event, err := a.GetEvent(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("%w; id=%s", ErrDocumentNotFound, id)
	}

	return a.storage.Update(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	event, err := a.GetEvent(ctx, id)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("%w; id=%s", ErrDocumentNotFound, id)
	}

	return a.storage.Delete(ctx, id)
}
