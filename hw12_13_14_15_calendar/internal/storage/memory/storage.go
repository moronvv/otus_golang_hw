package memorystorage

import (
	"context"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type InMemoryStorage struct {
	events map[int64]models.Event
	seqID  int64
}

func NewStorage(*config.DatabaseConf) storage.Storage {
	return &InMemoryStorage{}
}

func (s *InMemoryStorage) Connect(context.Context) error {
	s.seqID = 0
	s.events = map[int64]models.Event{}

	return nil
}

func (s *InMemoryStorage) Close(context.Context) error {
	return nil
}
