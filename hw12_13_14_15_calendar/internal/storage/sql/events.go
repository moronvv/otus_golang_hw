package sqlstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
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

type SqlEventStorage struct {
	store *SqlStorage
}

func NewEventStorage(store *SqlStorage) *SqlEventStorage {
	return &SqlEventStorage{
		store: store,
	}
}

func (s *SqlEventStorage) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	stmt, err := s.store.DB.PrepareNamedContext(ctx, createQuery)
	if err != nil {
		return nil, err
	}

	if err := stmt.Get(&event.ID, &event); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *SqlEventStorage) GetMany(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	if err := s.store.DB.SelectContext(ctx, &events, getManyQuery); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *SqlEventStorage) GetOne(ctx context.Context, id int64) (*models.Event, error) {
	var event models.Event
	if err := s.store.DB.GetContext(ctx, &event, getOneQuery, id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &event, nil
}

func (s *SqlEventStorage) Update(ctx context.Context, id int64, event *models.Event) (*models.Event, error) {
	event.ID = id
	if _, err := s.store.DB.NamedExecContext(ctx, updateQuery, event); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *SqlEventStorage) Delete(ctx context.Context, id int64) error {
	if _, err := s.store.DB.ExecContext(ctx, deleteQuery, id); err != nil {
		return err
	}

	return nil
}
