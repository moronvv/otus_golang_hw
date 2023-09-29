package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Datetime     time.Time
	Title        string
	Description  string
	Duration     time.Duration
	NotifyBefore time.Duration
	ID           uuid.UUID
	UserId       uuid.UUID
}
