package app_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	mockedstorage "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/mocked"
)

func getTestEvent() *models.Event {
	return &models.Event{
		Title: "test",
	}
}

func TestEventOperations(t *testing.T) {
	ctx := context.Background()
	mockStorage := mockedstorage.NewMockStorage(t)
	application := app.New(&slog.Logger{}, mockStorage)

	t.Run("create event", func(t *testing.T) {
		testEvent := getTestEvent()

		mockStorage.EXPECT().CreateEvent(mock.Anything, testEvent).Return(testEvent, nil).Once()

		event, err := application.CreateEvent(ctx, testEvent)
		require.NoError(t, err)
		require.Equal(t, testEvent, event)

		mockStorage.AssertExpectations(t)
	})

	t.Run("get events", func(t *testing.T) {
		testEvent := getTestEvent()

		mockStorage.EXPECT().GetEvents(mock.Anything).Return([]models.Event{*testEvent}, nil).Once()

		events, err := application.GetEvents(ctx)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, testEvent, &events[0])

		mockStorage.AssertExpectations(t)
	})

	t.Run("get event", func(t *testing.T) {
		id := uuid.New()
		testEvent := getTestEvent()

		mockStorage.EXPECT().GetEvent(mock.Anything, id).Return(testEvent, nil).Once()

		event, err := application.GetEvent(ctx, id)
		require.NoError(t, err)
		require.Equal(t, testEvent, event)

		mockStorage.AssertExpectations(t)
	})

	t.Run("get event no doc", func(t *testing.T) {
		id := uuid.New()

		mockStorage.EXPECT().GetEvent(mock.Anything, id).Return(nil, nil).Once()

		event, err := application.GetEvent(ctx, id)
		require.ErrorIs(t, err, app.ErrDocumentNotFound)
		require.Empty(t, event)

		mockStorage.AssertExpectations(t)
	})

	t.Run("update event", func(t *testing.T) {
		id := uuid.New()
		testEvent := getTestEvent()

		mockStorage.EXPECT().GetEvent(mock.Anything, id).Return(testEvent, nil).Once()
		mockStorage.EXPECT().UpdateEvent(mock.Anything, id, testEvent).Return(testEvent, nil).Once()

		event, err := application.UpdateEvent(ctx, id, testEvent)
		require.NoError(t, err)
		require.Equal(t, testEvent, event)

		mockStorage.AssertExpectations(t)
	})

	t.Run("update event no doc", func(t *testing.T) {
		id := uuid.New()
		testEvent := getTestEvent()

		mockStorage.EXPECT().GetEvent(mock.Anything, id).Return(nil, nil).Once()

		event, err := application.UpdateEvent(ctx, id, testEvent)
		require.ErrorIs(t, err, app.ErrDocumentNotFound)
		require.Empty(t, event)

		mockStorage.AssertExpectations(t)
	})

	t.Run("delete event", func(t *testing.T) {
		id := uuid.New()
		testEvent := getTestEvent()

		mockStorage.EXPECT().GetEvent(mock.Anything, id).Return(testEvent, nil).Once()
		mockStorage.EXPECT().DeleteEvent(mock.Anything, id).Return(nil)

		err := application.DeleteEvent(ctx, id)
		require.NoError(t, err)

		mockStorage.AssertExpectations(t)
	})

	t.Run("delete event no doc", func(t *testing.T) {
		id := uuid.New()

		mockStorage.EXPECT().GetEvent(mock.Anything, id).Return(nil, nil).Once()

		err := application.DeleteEvent(ctx, id)
		require.ErrorIs(t, err, app.ErrDocumentNotFound)

		mockStorage.AssertExpectations(t)
	})
}
