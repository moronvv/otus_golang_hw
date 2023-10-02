package memorystorage

import (
	"context"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
)

type InMemoryStorage struct {
	seqId  int64
	events map[int64]models.Event
}

func NewStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

func (s *InMemoryStorage) Connect(ctx context.Context) error {
	s.seqId = 0
	s.events = map[int64]models.Event{}

	return nil
}

func (s *InMemoryStorage) Close(ctx context.Context) error {
	return nil
}
