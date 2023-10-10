package app_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
	mockedstorage "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/mocked"
)

func getTestEvent() *models.Event {
	return &models.Event{
		Title: "test",
	}
}

func TestEventOperations(t *testing.T) {
	ctx := context.Background()
	mockEventStorage := mockedstorage.NewMockEventStorage(t)
	application := app.New(&slog.Logger{}, &storage.Storages{
		Events: mockEventStorage,
	})

	t.Run("create event", func(t *testing.T) {
		testEvent := getTestEvent()

		mockEventStorage.EXPECT().Create(mock.Anything, testEvent).Return(testEvent, nil).Once()

		event, err := application.CreateEvent(ctx, testEvent)
		require.NoError(t, err)
		require.Equal(t, testEvent, event)

		mockEventStorage.AssertExpectations(t)
	})

	t.Run("get events", func(t *testing.T) {
		testEvent := getTestEvent()

		mockEventStorage.EXPECT().GetMany(mock.Anything).Return([]models.Event{*testEvent}, nil).Once()

		events, err := application.GetEvents(ctx)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, testEvent, &events[0])

		mockEventStorage.AssertExpectations(t)
	})

	t.Run("get event", func(t *testing.T) {
		var id int64 = 1
		testEvent := getTestEvent()

		mockEventStorage.EXPECT().GetOne(mock.Anything, id).Return(testEvent, nil).Once()

		event, err := application.GetEvent(ctx, id)
		require.NoError(t, err)
		require.Equal(t, testEvent, event)

		mockEventStorage.AssertExpectations(t)
	})

	t.Run("get event no doc", func(t *testing.T) {
		var id int64 = 1

		mockEventStorage.EXPECT().GetOne(mock.Anything, id).Return(nil, nil).Once()

		event, err := application.GetEvent(ctx, id)
		require.ErrorIs(t, err, app.ErrDocumentNotFound)
		require.Empty(t, event)

		mockEventStorage.AssertExpectations(t)
	})

	t.Run("update event", func(t *testing.T) {
		var id int64 = 1
		testEvent := getTestEvent()

		mockEventStorage.EXPECT().Update(mock.Anything, id, testEvent).Return(testEvent, nil).Once()

		event, err := application.UpdateEvent(ctx, id, testEvent)
		require.NoError(t, err)
		require.Equal(t, testEvent, event)

		mockEventStorage.AssertExpectations(t)
	})

	t.Run("update event no doc", func(t *testing.T) {
		var id int64 = 1
		testEvent := getTestEvent()

		mockEventStorage.EXPECT().Update(mock.Anything, id, testEvent).Return(nil, nil).Once()

		event, err := application.UpdateEvent(ctx, id, testEvent)
		require.ErrorIs(t, err, app.ErrDocumentNotFound)
		require.Empty(t, event)

		mockEventStorage.AssertExpectations(t)
	})

	t.Run("delete event", func(t *testing.T) {
		var id int64 = 1

		mockEventStorage.EXPECT().Delete(mock.Anything, id).Return(true, nil).Once()

		err := application.DeleteEvent(ctx, id)
		require.NoError(t, err)

		mockEventStorage.AssertExpectations(t)
	})

	t.Run("delete event no doc", func(t *testing.T) {
		var id int64 = 1

		mockEventStorage.EXPECT().Delete(mock.Anything, id).Return(false, nil).Once()

		err := application.DeleteEvent(ctx, id)
		require.ErrorIs(t, err, app.ErrDocumentNotFound)

		mockEventStorage.AssertExpectations(t)
	})
}
