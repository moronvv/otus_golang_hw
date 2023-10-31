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
	internalcontext "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/context"
	internalerrors "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/errors"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	sqlstorage "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

func getTestEvent() *models.Event {
	return &models.Event{
		Title: "test",
		Description: sql.NullString{
			String: "description",
			Valid:  true,
		},
		Datetime: time.Now().UTC(),
		Duration: 1 * time.Minute,
		UserID:   uuid.New(),
	}
}

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
	userID := uuid.New()
	ctx := internalcontext.SetUserID(context.Background(), userID)
	eventStore := s.eventStore

	// create foreign event
	foreignEvent := getTestEvent()
	_, err := eventStore.Create(ctx, foreignEvent)
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
		Title: "updated",
		Description: sql.NullString{
			String: "description",
			Valid:  true,
		},
		Datetime: time.Now().UTC(),
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

func (s *SQLStorageSuite) TestSqlStorageDocOperationForbidden() {
	t := s.T()
	ctx := internalcontext.SetUserID(context.Background(), uuid.New())
	eventStore := s.eventStore

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

func (s *SQLStorageSuite) TestSqlStorageDocNotFound() {
	t := s.T()
	ctx := internalcontext.SetUserID(context.Background(), uuid.New())
	eventStore := s.eventStore
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

func TestSqlStorageSuite(t *testing.T) {
	t.Skip("Skip before setting postgres container in ci/cd")

	suite.Run(t, new(SQLStorageSuite))
}
