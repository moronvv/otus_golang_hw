package internalhttproutes

import (
	"net/http"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
)

func getEvents(app *app.App) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		// events, err := app.GetEvents(ctx)
	}
}

func getEvent(app *app.App) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func createEvent(app *app.App) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func updateEvent(app *app.App) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func deleteEvent(app *app.App) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {}
}
