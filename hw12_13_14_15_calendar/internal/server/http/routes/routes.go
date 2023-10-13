package internalhttproutes

import (
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
)

type components struct {
	App       app.App
	Validator *validator.Validate
}

type handlerFn func(http.ResponseWriter, *http.Request)

func SetupRoutes(app app.App) *mux.Router {
	cmps := &components{
		App:       app,
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	router := mux.NewRouter()

	// events
	eventsRouter := router.PathPrefix("/events").Subrouter()
	eventsRouter.HandleFunc("", getEvents(cmps)).Methods("GET")
	eventsRouter.HandleFunc("/{id:[0-9]+}", getEvent(cmps)).Methods("GET")
	eventsRouter.HandleFunc("", createEvent(cmps)).Methods("POST")
	eventsRouter.HandleFunc("/{id:[0-9]+}", updateEvent(cmps)).Methods("PUT")
	eventsRouter.HandleFunc("/{id:[0-9]+}", deleteEvent(cmps)).Methods("DELETE")

	return router
}
