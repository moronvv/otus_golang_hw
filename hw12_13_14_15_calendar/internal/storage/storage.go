package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
)

type Storage interface {
	CreateEvent(context.Context, *models.Event) (*models.Event, error)
	GetEvents(context.Context) ([]models.Event, error)
	GetEvent(context.Context, uuid.UUID) (*models.Event, error)
	UpdateEvent(context.Context, uuid.UUID, *models.Event) (*models.Event, error)
	DeleteEvent(context.Context, uuid.UUID) error
}
