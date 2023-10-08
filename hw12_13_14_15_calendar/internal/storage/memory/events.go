package memorystorage

import (
	"context"
	"sync"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type InMemoryEventStorage struct {
	store *InMemoryStorage
	mu    sync.RWMutex
}

func NewEventStorage(store storage.Storage) storage.EventStorage {
	return &InMemoryEventStorage{
		store: store.(*InMemoryStorage),
		mu:    sync.RWMutex{},
	}
}

func (s *InMemoryEventStorage) Create(_ context.Context, event *models.Event) (*models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store.seqID++
	event.ID = s.store.seqID
	s.store.events[event.ID] = *event

	return event, nil
}

func (s *InMemoryEventStorage) GetMany(_ context.Context) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := []models.Event{}
	for _, event := range s.store.events {
		events = append(events, event)
	}

	return events, nil
}

func (s *InMemoryEventStorage) GetOne(_ context.Context, id int64) (*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.store.events[id]
	if !ok {
		return nil, nil
	}

	return &event, nil
}

func (s *InMemoryEventStorage) Update(_ context.Context, id int64, event *models.Event) (*models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = id
	s.store.events[event.ID] = *event

	return event, nil
}

func (s *InMemoryEventStorage) Delete(_ context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store.events, id)

	return nil
}
