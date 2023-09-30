package sqlstorage

import (
	"context"

	"github.com/google/uuid"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
)

type SqlEventStorage struct {
	store *SqlStorage
}

func NewEventStorage(store *SqlStorage) *SqlEventStorage {
	return &SqlEventStorage{
		store: store,
	}
}

func (s *SqlEventStorage) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	// TODO: implement
	return nil, nil
}

func (s *SqlEventStorage) GetMany(ctx context.Context) ([]models.Event, error) {
	// TODO: implement
	return nil, nil
}

func (s *SqlEventStorage) GetOne(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	// TODO: implement
	return nil, nil
}

func (s *SqlEventStorage) Update(ctx context.Context, id uuid.UUID, event *models.Event) (*models.Event, error) {
	// TODO: implement
	return nil, nil
}

func (s *SqlEventStorage) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: implement
	return nil
}
