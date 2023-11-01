package memorystorage

import (
	"context"
	"sync"

	internalcontext "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/context"
	internalerrors "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/errors"
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

func (s *InMemoryEventStorage) GetMany(ctx context.Context) ([]models.Event, error) {
	userID := internalcontext.GetUserID(ctx)

	s.mu.RLock()
	defer s.mu.RUnlock()

	events := []models.Event{}
	for _, event := range s.store.events {
		if event.UserID == userID {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *InMemoryEventStorage) GetOne(ctx context.Context, id int64) (*models.Event, error) {
	userID := internalcontext.GetUserID(ctx)

	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.store.events[id]
	if !ok {
		return nil, internalerrors.ErrDocumentNotFound
	}
	if event.UserID != userID {
		return nil, internalerrors.ErrDocumentOperationForbidden
	}

	return &event, nil
}

func (s *InMemoryEventStorage) Update(ctx context.Context, id int64, updEvent *models.Event) (*models.Event, error) {
	userID := internalcontext.GetUserID(ctx)

	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.store.events[id]
	if !ok {
		return nil, internalerrors.ErrDocumentNotFound
	}
	if event.UserID != userID {
		return nil, internalerrors.ErrDocumentOperationForbidden
	}

	updEvent.ID = id
	s.store.events[updEvent.ID] = *updEvent

	return updEvent, nil
}

func (s *InMemoryEventStorage) Delete(ctx context.Context, id int64) error {
	userID := internalcontext.GetUserID(ctx)

	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.store.events[id]
	if !ok {
		return internalerrors.ErrDocumentNotFound
	}
	if event.UserID != userID {
		return internalerrors.ErrDocumentOperationForbidden
	}

	delete(s.store.events, id)
	return nil
}
