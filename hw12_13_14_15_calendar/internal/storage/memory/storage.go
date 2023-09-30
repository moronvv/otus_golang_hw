package memorystorage

import (
	"context"

	"github.com/google/uuid"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
)

type InMemoryStorage struct {
	events map[uuid.UUID]models.Event
}

func NewStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

func (s *InMemoryStorage) Connect(ctx context.Context) error {
	s.events = map[uuid.UUID]models.Event{}

	return nil
}

func (s *InMemoryStorage) Close(ctx context.Context) error {
	return nil
}
