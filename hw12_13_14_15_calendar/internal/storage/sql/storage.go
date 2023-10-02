package sqlstorage

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib" // postgres driver
)

type SQLStorage struct {
	DB  *sqlx.DB
	cfg *config.DatabaseConf
}

func NewStorage(cfg *config.DatabaseConf) *SQLStorage {
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

	s.DB = db
	return nil
}

func (s *SQLStorage) Close(context.Context) error {
	return s.DB.Close()
}
