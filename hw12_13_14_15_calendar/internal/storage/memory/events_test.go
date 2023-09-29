package memorystorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	memorystorage "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/memory"
)

func TestEventStorage(t *testing.T) {
	eventStore := memorystorage.New()
	ctx := context.Background()

	testEvent := &models.Event{
		Title:    "test",
		Datetime: time.Now(),
		Duration: 1 * time.Minute,
	}

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
	}
	event, err = eventStore.Update(ctx, id, updatedTestEvent)
	require.NoError(t, err)
	require.Equal(t, updatedTestEvent, event)

	// delete
	err = eventStore.Delete(ctx, id)
	require.NoError(t, err)
	event, err = eventStore.GetOne(ctx, id)
	require.NoError(t, err)
	require.Empty(t, event)
}
