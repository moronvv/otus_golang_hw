package memorystorage

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	store map[uuid.UUID]storage.Event
	mu    sync.RWMutex
}

func New() *Storage {
	return &Storage{
		store: map[uuid.UUID]storage.Event{},
		mu:    sync.RWMutex{},
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) (*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = uuid.New()
	s.store[event.ID] = *event

	return event, nil
}

func (s *Storage) GetEvents(ctx context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []storage.Event
	for _, event := range s.store {
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) GetEvent(ctx context.Context, id uuid.UUID) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.store[id]
	if !ok {
		return nil, nil
	}

	return &event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id uuid.UUID, event *storage.Event) (*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = id
	s.store[event.ID] = *event

	return event, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, id)

	return nil
}
