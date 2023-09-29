package storage

import (
	"context"

	"github.com/google/uuid"
)

type Storage interface {
	CreateEvent(context.Context, *Event) (*Event, error)
	GetEvents(context.Context) ([]Event, error)
	GetEvent(context.Context, uuid.UUID) (*Event, error)
	UpdateEvent(context.Context, uuid.UUID, *Event) (*Event, error)
	DeleteEvent(context.Context, uuid.UUID) error
}
