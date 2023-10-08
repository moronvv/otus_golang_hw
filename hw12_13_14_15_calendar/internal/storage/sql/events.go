package sqlstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
)

const (
	createQuery = `
    INSERT INTO events (title, description, datetime, duration, user_id, notify_before)
    VALUES (:title, :description, :datetime, :duration, :user_id, :notify_before)
    RETURNING id;
  `
	getManyQuery = `
    SELECT * FROM events;
  `
	getOneQuery = `
    SELECT * FROM events
    WHERE id = $1;
  `
	updateQuery = `
    UPDATE events
    SET
      title = :title,
      description = :description,
      datetime = :datetime,
      user_id = :user_id,
      notify_before = :notify_before
    WHERE id = :id;
  `
	deleteQuery = `
    DELETE FROM events
    WHERE id = $1;
  `
)

type SQLEventStorage struct {
	store *SQLStorage
}

func NewEventStorage(store storage.Storage) storage.EventStorage {
	return &SQLEventStorage{
		store: store.(*SQLStorage),
	}
}

func (s *SQLEventStorage) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	stmt, err := s.store.db.PrepareNamedContext(ctx, createQuery)
	if err != nil {
		return nil, err
	}

	if err := stmt.Get(&event.ID, &event); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *SQLEventStorage) GetMany(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	if err := s.store.db.SelectContext(ctx, &events, getManyQuery); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *SQLEventStorage) GetOne(ctx context.Context, id int64) (*models.Event, error) {
	var event models.Event
	if err := s.store.db.GetContext(ctx, &event, getOneQuery, id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &event, nil
}

func (s *SQLEventStorage) Update(ctx context.Context, id int64, event *models.Event) (*models.Event, error) {
	event.ID = id
	if _, err := s.store.db.NamedExecContext(ctx, updateQuery, event); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *SQLEventStorage) Delete(ctx context.Context, id int64) error {
	if _, err := s.store.db.ExecContext(ctx, deleteQuery, id); err != nil {
		return err
	}

	return nil
}
