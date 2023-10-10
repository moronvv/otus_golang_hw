package sqlstorage

import (
	"context"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib" // postgres driver
	"github.com/jmoiron/sqlx"
	goose "github.com/pressly/goose/v3"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/utils"
)

var migrationsDir = filepath.Join(utils.Root, "migrations")

type SQLStorage struct {
	db  *sqlx.DB
	cfg *config.DatabaseConf
}

func NewStorage(cfg *config.DatabaseConf) storage.Storage {
	return &SQLStorage{
		cfg: cfg,
	}
}

func (s *SQLStorage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "pgx", s.cfg.DSN)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(s.cfg.MaxOpenConns)
	db.SetMaxIdleConns(s.cfg.MaxIdleConns)
	db.SetConnMaxLifetime(s.cfg.ConnMaxLifetime)

	s.db = db
	return nil
}

func (s *SQLStorage) Close(context.Context) error {
	return s.db.Close()
}

func (s *SQLStorage) MigrateAllUp(ctx context.Context) error {
	return goose.UpContext(ctx, s.db.DB, migrationsDir)
}

func (s *SQLStorage) MigrateAllDown(ctx context.Context) error {
	return goose.DownToContext(ctx, s.db.DB, migrationsDir, 0)
}
