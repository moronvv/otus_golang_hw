package sqlstorage_test

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	sqlstorage "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type SqlStorageSuite struct {
	suite.Suite
	store      *sqlstorage.SqlStorage
	eventStore *sqlstorage.SqlEventStorage
}

func (s *SqlStorageSuite) SetupSuite() {
	cfg := &config.DatabaseConf{
		DSN:             "postgresql://",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: 30 * time.Second,
	}
	s.store = sqlstorage.NewStorage(cfg)

	err := s.store.Connect(context.Background())
	require.NoError(s.T(), err)

	// TODO: run migrations up

	s.eventStore = sqlstorage.NewEventStorage(s.store)
}

func (s *SqlStorageSuite) TearDownSuite() {
	// TODO: run migrations down

	err := s.store.Close(context.Background())
	require.NoError(s.T(), err)
}

func (s *SqlStorageSuite) TestEventStorage() {
}

// func TestSqlStorageSuite(t *testing.T) {
// }
