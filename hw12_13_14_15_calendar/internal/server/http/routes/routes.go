package internalhttproutes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
)

type handlerFn func(http.ResponseWriter, *http.Request)

func SetupRoutes(app *app.App) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", pingHandler).Methods("GET")

	// events
	eventsRouter := router.PathPrefix("/events").Subrouter()
	eventsRouter.HandleFunc("", getEvents(app)).Methods("GET")
	eventsRouter.HandleFunc("/{id}", getEvent(app)).Methods("GET")
	eventsRouter.HandleFunc("", createEvent(app)).Methods("POST")
	eventsRouter.HandleFunc("/{id}", updateEvent(app)).Methods("PUT")
	eventsRouter.HandleFunc("/{id}", deleteEvent(app)).Methods("DELETE")

	return router
}
