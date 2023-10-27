package storage

import (
	"context"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
)

type Storage interface {
	Connect(context.Context) error
	Close(context.Context) error
}

type EventStorage interface {
	Create(context.Context, *models.Event) (*models.Event, error)
	GetMany(context.Context) ([]models.Event, error)
	GetOne(context.Context, int64) (*models.Event, error)
	Update(context.Context, int64, *models.Event) (*models.Event, error)
	Delete(context.Context, int64) error
}

type Storages struct {
	DB     Storage
	Events EventStorage
}
