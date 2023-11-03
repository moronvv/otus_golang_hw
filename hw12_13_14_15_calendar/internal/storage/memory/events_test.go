package memorystorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	internalcontext "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/context"
	internalerrors "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/errors"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/memory"
)

func getTestEvent() *models.Event {
	return &models.Event{
		Title:    "test",
		Datetime: time.Now(),
		Duration: 1 * time.Minute,
		UserID:   uuid.New(),
	}
}

func getEventStorage(ctx context.Context) (storage.EventStorage, error) {
	store := memorystorage.NewStorage(nil)
	if err := store.Connect(ctx); err != nil {
		return nil, err
	}

	return memorystorage.NewEventStorage(store), nil
}

func TestEventStorage(t *testing.T) {
	userID := uuid.New()
	ctx := internalcontext.SetUserID(context.Background(), userID)
	eventStore, err := getEventStorage(ctx)
	require.NoError(t, err)

	// create foreign event
	foreignEvent := getTestEvent()
	_, err = eventStore.Create(ctx, foreignEvent)
	require.NoError(t, err)

	testEvent := getTestEvent()
	testEvent.UserID = userID

	// create
	event, err := eventStore.Create(ctx, testEvent)
	require.NoError(t, err)
	require.Equal(t, testEvent, event)
	id := event.ID

	// read
	events, err := eventStore.GetMany(ctx)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, testEvent, &events[0])
	event, err = eventStore.GetOne(ctx, id)
	require.NoError(t, err)
	require.Equal(t, testEvent, event)

	// update
	updatedTestEvent := &models.Event{
		Title:    "updated",
		Datetime: time.Now(),
		Duration: 2 * time.Minute,
		UserID:   userID,
	}
	event, err = eventStore.Update(ctx, id, updatedTestEvent)
	require.NoError(t, err)
	require.Equal(t, updatedTestEvent, event)

	// delete
	err = eventStore.Delete(ctx, id)
	require.NoError(t, err)
	event, err = eventStore.GetOne(ctx, id)
	require.ErrorIs(t, err, internalerrors.ErrDocumentNotFound)
	require.Empty(t, event)
}

func TestEventStorageDocOperationForbidden(t *testing.T) {
	ctx := internalcontext.SetUserID(context.Background(), uuid.New())
	eventStore, err := getEventStorage(ctx)
	require.NoError(t, err)

	testEvent := getTestEvent()
	event, err := eventStore.Create(ctx, testEvent)
	require.NoError(t, err)
	id := event.ID

	// read
	event, err = eventStore.GetOne(ctx, id)
	require.ErrorIs(t, err, internalerrors.ErrDocumentOperationForbidden)
	require.Nil(t, event)

	// update
	event, err = eventStore.Update(ctx, id, &models.Event{Title: "updated"})
	require.ErrorIs(t, err, internalerrors.ErrDocumentOperationForbidden)
	require.Nil(t, event)

	// delete
	err = eventStore.Delete(ctx, id)
	require.ErrorIs(t, err, internalerrors.ErrDocumentOperationForbidden)
}

func TestEventStorageDocNotFound(t *testing.T) {
	userID := uuid.New()
	ctx := internalcontext.SetUserID(context.Background(), userID)
	eventStore, err := getEventStorage(ctx)
	require.NoError(t, err)
	var id int64 = 1337

	// read
	event, err := eventStore.GetOne(ctx, id)
	require.ErrorIs(t, err, internalerrors.ErrDocumentNotFound)
	require.Nil(t, event)

	// update
	event, err = eventStore.Update(ctx, id, &models.Event{Title: "updated"})
	require.ErrorIs(t, err, internalerrors.ErrDocumentNotFound)
	require.Nil(t, event)

	// delete
	err = eventStore.Delete(ctx, id)
	require.ErrorIs(t, err, internalerrors.ErrDocumentNotFound)
}
