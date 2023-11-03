package internalhttproutes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
)

type handlerFn func(http.ResponseWriter, *http.Request)

func SetupRoutes(app app.App) *mux.Router {
	router := mux.NewRouter()

	// events
	eventsRouter := router.PathPrefix("/events").Subrouter()
	eventsRouter.HandleFunc("", createEvent(app)).Methods("POST")
	eventsRouter.HandleFunc("", getEvents(app)).Methods("GET")
	eventsRouter.HandleFunc("/{id:[0-9]+}", getEvent(app)).Methods("GET")
	eventsRouter.HandleFunc("/{id:[0-9]+}", updateEvent(app)).Methods("PUT")
	eventsRouter.HandleFunc("/{id:[0-9]+}", deleteEvent(app)).Methods("DELETE")

	return router
}
