package internalhttproutes

import (
	"net/http"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
)

func getEvents(app app.App) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := app.GetEvents(r.Context())
		if err != nil {
			setErrorReponse(w, http.StatusInternalServerError, err)
			return
		}

		setJSONResponse(w, http.StatusOK, events)
	}
}

func getEvent(app app.App) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			setErrorReponse(w, http.StatusBadRequest, err)
			return
		}

		event, err := app.GetEvent(r.Context(), id)
		if err != nil {
			setErrorReponse(w, getErrorStatus(err), err)
			return
		}

		setJSONResponse(w, http.StatusOK, event)
	}
}

func createEvent(app app.App) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := getPayload[models.Event](r)
		if err != nil {
			setErrorReponse(w, http.StatusBadRequest, err)
			return
		}

		createdEvent, err := app.CreateEvent(r.Context(), request)
		if err != nil {
			setErrorReponse(w, getErrorStatus(err), err)
			return
		}

		setJSONResponse(w, http.StatusCreated, createdEvent)
	}
}

func updateEvent(app app.App) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			setErrorReponse(w, http.StatusBadRequest, err)
			return
		}

		request, err := getPayload[models.Event](r)
		if err != nil {
			setErrorReponse(w, http.StatusBadRequest, err)
			return
		}

		updatedEvent, err := app.UpdateEvent(r.Context(), id, request)
		if err != nil {
			setErrorReponse(w, getErrorStatus(err), err)
			return
		}

		setJSONResponse(w, http.StatusOK, updatedEvent)
	}
}

func deleteEvent(app app.App) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			setErrorReponse(w, http.StatusBadRequest, err)
			return
		}

		if err := app.DeleteEvent(r.Context(), id); err != nil {
			setErrorReponse(w, getErrorStatus(err), err)
			return
		}

		setEmptyResponse(w, http.StatusNoContent)
	}
}
