package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Datetime     time.Time      `db:"datetime"`
	Title        string         `db:"title"`
	Description  sql.NullString `db:"description"`
	ID           int64          `db:"id"`
	NotifyBefore time.Duration  `db:"notify_before"`
	Duration     time.Duration  `db:"duration"`
	UserId       uuid.UUID      `db:"user_id"`
}
