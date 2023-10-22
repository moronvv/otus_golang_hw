package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/pb"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	requestValidator  *validator.Validate
)

func init() {
	requestValidator = validator.New(validator.WithRequiredStructEnabled())
}

type Event struct {
	ID           int64          `db:"id"`
	Title        string         `db:"title"`
	Description  sql.NullString `db:"description"`
	Datetime     time.Time      `db:"datetime"`
	Duration     time.Duration  `db:"duration"`
	UserID       uuid.UUID      `db:"user_id"`
	NotifyBefore time.Duration  `db:"notify_before"`
}

type eventRequest struct {
	Title        string `json:"title" validate:"required,min=2,max=30"`
	Description  string `json:"description" validate:"omitempty,min=5,max=150"`
	Datetime     string `json:"datetime" validate:"required"`
	Duration     string `json:"duration" validate:"required"`
	UserID       string `json:"user_id" validate:"required,uuid4"`
	NotifyBefore string `json:"notify_before"`
}

type eventResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	Datetime     time.Time `json:"datetime"`
	Duration     string    `json:"duration"`
	UserID       uuid.UUID `json:"user_id"`
	NotifyBefore string    `json:"notify_before,omitempty"`
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var err error

	var req eventRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	if err := requestValidator.Struct(req); err != nil {
		return fmt.Errorf("%w; %w", ErrInvalidRequest, err)
	}

	e.Title = req.Title

	// description
	if req.Description != "" {
		e.Description = sql.NullString{
			String: req.Description,
			Valid:  true,
		}
	}

	// datetime
	e.Datetime, err = time.Parse(time.RFC3339Nano, req.Datetime)
	if err != nil {
		return fmt.Errorf("datetime parse error; %w", err)
	}

	// duration
	e.Duration, err = time.ParseDuration(req.Duration)
	if err != nil {
		return fmt.Errorf("duration parse error; %w", err)
	}

	// user_id
	e.UserID, err = uuid.Parse(req.UserID)
	if err != nil {
		return fmt.Errorf("user_id parse error; %w", err)
	}

	// notify_before
	if req.NotifyBefore != "" {
		e.NotifyBefore, err = time.ParseDuration(req.NotifyBefore)
		if err != nil {
			return fmt.Errorf("notify_before parse error; %w", err)
		}
	}

	return nil
}

func (e *Event) MarshalJSON() ([]byte, error) {
	// description
	description := ""
	if e.Description.Valid {
		description = e.Description.String
	}

	// notifyBefore
	notifyBefore := ""
	if e.NotifyBefore != 0 {
		notifyBefore = e.NotifyBefore.String()
	}

	resp := eventResponse{
		ID:           e.ID,
		Title:        e.Title,
		Description:  description,
		Datetime:     e.Datetime,
		Duration:     e.Duration.String(),
		UserID:       e.UserID,
		NotifyBefore: notifyBefore,
	}

	return json.Marshal(&resp)
}

func (e *Event) UnmarshalPB(req *pb.EventRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	var err error

	e.ID = 0
	e.Title = req.Title
	if req.Description != "" {
		e.Description = sql.NullString{
			String: req.Description,
			Valid:  true,
		}
	}
	e.Datetime = req.Datetime.AsTime()
	e.Duration = req.Duration.AsDuration()
	e.UserID, err = uuid.Parse(req.UserId)
	if err != nil {
		return err
	}
	e.NotifyBefore = req.NotifyBefore.AsDuration()

	return nil
}

func (e *Event) MarshalPB() *pb.EventResponse {
	return &pb.EventResponse{
		Id:           e.ID,
		Title:        e.Title,
		Description:  e.Description.String,
		Datetime:     timestamppb.New(e.Datetime),
		Duration:     durationpb.New(e.Duration),
		UserId:       e.UserID.String(),
		NotifyBefore: durationpb.New(e.NotifyBefore),
	}
}
