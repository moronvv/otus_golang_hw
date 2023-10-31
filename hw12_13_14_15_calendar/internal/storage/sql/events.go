package sqlstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	internalcontext "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/context"
	internalerrors "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/errors"
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
	userID := internalcontext.GetUserID(ctx)
	var events []models.Event

	query := `
    SELECT * FROM events
    WHERE user_id = $1;
  `
	if err := s.store.db.SelectContext(ctx, &events, query, userID); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *SQLEventStorage) GetOne(ctx context.Context, id int64) (*models.Event, error) {
	userID := internalcontext.GetUserID(ctx)
	var event models.Event

	query := `
    SELECT * FROM events
    WHERE id = $1;
  `
	if err := s.store.db.GetContext(ctx, &event, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalerrors.ErrDocumentNotFound
		}

		return nil, err
	}
	if event.UserID != userID {
		return nil, internalerrors.ErrDocumentOperationForbidden
	}

	return &event, nil
}

func (s *SQLEventStorage) Update(ctx context.Context, id int64, updEvent *models.Event) (*models.Event, error) {
	userID := internalcontext.GetUserID(ctx)

	tx, err := s.store.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
    SELECT user_id FROM events
    WHERE id = $1
    FOR UPDATE;
  `
	var eventUserID uuid.UUID
	if err := tx.GetContext(ctx, &eventUserID, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalerrors.ErrDocumentNotFound
		}

		return nil, err
	}
	if eventUserID != userID {
		return nil, internalerrors.ErrDocumentOperationForbidden
	}

	updEvent.ID = id
	query = `
    UPDATE events
    SET
      title = :title,
      description = :description,
      datetime = :datetime,
      user_id = :user_id,
      notify_before = :notify_before
    WHERE id = :id;
  `
	_, err = tx.NamedExecContext(ctx, query, updEvent)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updEvent, nil
}

func (s *SQLEventStorage) Delete(ctx context.Context, id int64) error {
	userID := internalcontext.GetUserID(ctx)

	tx, err := s.store.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
    SELECT user_id FROM events
    WHERE id = $1
    FOR UPDATE;
  `
	var eventUserID uuid.UUID
	if err := tx.GetContext(ctx, &eventUserID, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internalerrors.ErrDocumentNotFound
		}

		return err
	}
	if eventUserID != userID {
		return internalerrors.ErrDocumentOperationForbidden
	}

	query = `
    DELETE FROM events
    WHERE id = $1;
  `
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
