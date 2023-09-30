package memorystorage

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
)

type InMemoryEventStorage struct {
	store *InMemoryStorage
	mu    sync.RWMutex
}

func NewEventStorage(store *InMemoryStorage) *InMemoryEventStorage {
	return &InMemoryEventStorage{
		store: store,
		mu:    sync.RWMutex{},
	}
}

func (s *InMemoryEventStorage) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = uuid.New()
	s.store.events[event.ID] = *event

	return event, nil
}

func (s *InMemoryEventStorage) GetMany(ctx context.Context) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []models.Event
	for _, event := range s.store.events {
		events = append(events, event)
	}

	return events, nil
}

func (s *InMemoryEventStorage) GetOne(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.store.events[id]
	if !ok {
		return nil, nil
	}

	return &event, nil
}

func (s *InMemoryEventStorage) Update(ctx context.Context, id uuid.UUID, event *models.Event) (*models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = id
	s.store.events[event.ID] = *event

	return event, nil
}

func (s *InMemoryEventStorage) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store.events, id)

	return nil
}
