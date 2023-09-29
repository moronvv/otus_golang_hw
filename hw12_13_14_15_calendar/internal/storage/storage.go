package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
)

type EventStorage interface {
	Create(context.Context, *models.Event) (*models.Event, error)
	GetMany(context.Context) ([]models.Event, error)
	GetOne(context.Context, uuid.UUID) (*models.Event, error)
	Update(context.Context, uuid.UUID, *models.Event) (*models.Event, error)
	Delete(context.Context, uuid.UUID) error
}
