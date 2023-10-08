package sqlstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
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
	query := `
    INSERT INTO events (title, description, datetime, duration, user_id, notify_before)
    VALUES (:title, :description, :datetime, :duration, :user_id, :notify_before)
    RETURNING id;
  `
	stmt, err := s.store.db.PrepareNamedContext(ctx, query)
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

	query := "SELECT * FROM events;"
	if err := s.store.db.SelectContext(ctx, &events, query); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *SQLEventStorage) GetOne(ctx context.Context, id int64) (*models.Event, error) {
	var event models.Event

	query := `
    SELECT * FROM events
    WHERE id = $1;
  `
	if err := s.store.db.GetContext(ctx, &event, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &event, nil
}

func (s *SQLEventStorage) Update(ctx context.Context, id int64, event *models.Event) (*models.Event, error) {
	event.ID = id

	query := `
    UPDATE events
    SET
      title = :title,
      description = :description,
      datetime = :datetime,
      user_id = :user_id,
      notify_before = :notify_before
    WHERE id = :id;
  `
	result, err := s.store.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return nil, err
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		return nil, err
	} else if rowsAffected == 0 {
		return nil, nil
	}

	return event, nil
}

func (s *SQLEventStorage) Delete(ctx context.Context, id int64) (bool, error) {
	query := `
    DELETE FROM events
    WHERE id = $1;
  `
	result, err := s.store.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else if rowsAffected == 0 {
		return false, nil
	}

	return true, err
}
