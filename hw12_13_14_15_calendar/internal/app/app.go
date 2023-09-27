package app

import (
	"context"
	"log/slog"
)

type App struct { // TODO
}

type Storage interface { // TODO
}

func New(logger *slog.Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
