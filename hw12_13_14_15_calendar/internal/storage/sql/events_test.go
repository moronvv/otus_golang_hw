package sqlstorage_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	sqlstorage "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type SQLStorageSuite struct {
	suite.Suite
	store      *sqlstorage.SQLStorage
	eventStore *sqlstorage.SQLEventStorage
}

func (s *SQLStorageSuite) SetupSuite() {
	ctx := context.Background()

	cfg := &config.DatabaseConf{
		DSN:             os.Getenv("CALENDAR_DATABASE_DSN"),
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: 30 * time.Second,
	}
	s.store = sqlstorage.NewStorage(cfg).(*sqlstorage.SQLStorage)

	err := s.store.Connect(context.Background())
	require.NoError(s.T(), err)

	err = s.store.MigrateAllUp(ctx)
	require.NoError(s.T(), err)

	s.eventStore = sqlstorage.NewEventStorage(s.store).(*sqlstorage.SQLEventStorage)
}

func (s *SQLStorageSuite) TearDownSuite() {
	ctx := context.Background()

	err := s.store.MigrateAllDown(ctx)
	require.NoError(s.T(), err)

	err = s.store.Close(context.Background())
	require.NoError(s.T(), err)
}

func (s *SQLStorageSuite) TestEventStorage() {
	t := s.T()
	ctx := context.Background()
	eventStore := s.eventStore

	testEvent := &models.Event{
		Title: "test",
		Description: sql.NullString{
			String: "description",
			Valid:  true,
		},
		Datetime: time.Now().UTC(),
		Duration: 1 * time.Minute,
		UserID:   uuid.New(),
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
		Title: "updated",
		Description: sql.NullString{
			String: "description",
			Valid:  true,
		},
		Datetime: time.Now().UTC(),
		Duration: 2 * time.Minute,
		UserID:   uuid.New(),
	}
	event, err = eventStore.Update(ctx, id, updatedTestEvent)
	require.NoError(t, err)
	require.Equal(t, updatedTestEvent, event)

	// delete
	ok, err := eventStore.Delete(ctx, id)
	require.NoError(t, err)
	require.True(t, ok)
	event, err = eventStore.GetOne(ctx, id)
	require.NoError(t, err)
	require.Empty(t, event)
}

func (s *SQLStorageSuite) TestSqlStorageSuiteDocNotFound() {
	t := s.T()
	ctx := context.Background()
	eventStore := s.eventStore
	var id int64 = 1337

	// read
	event, err := eventStore.GetOne(ctx, id)
	require.NoError(t, err)
	require.Nil(t, event)

	// update
	event, err = eventStore.Update(ctx, id, &models.Event{Title: "updated"})
	require.NoError(t, err)
	require.Nil(t, event)

	// delete
	ok, err := eventStore.Delete(ctx, id)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestSqlStorageSuite(t *testing.T) {
	t.Skip("Skip before setting postgres container in ci/cd")

	suite.Run(t, new(SQLStorageSuite))
}
