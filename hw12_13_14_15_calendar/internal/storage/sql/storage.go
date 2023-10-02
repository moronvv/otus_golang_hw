package sqlstorage

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
)

type SqlStorage struct {
	DB  *sqlx.DB
	cfg *config.DatabaseConf
}

func NewStorage(cfg *config.DatabaseConf) *SqlStorage {
	return &SqlStorage{
		cfg: cfg,
	}
}

func (s *SqlStorage) Connect(ctx context.Context) error {
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

func (s *SqlStorage) Close(ctx context.Context) error {
	return s.DB.Close()
}
