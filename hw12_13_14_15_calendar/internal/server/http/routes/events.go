package internalhttproutes

import (
	"net/http"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
)

func getEvents(cmps *components) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := cmps.App.GetEvents(r.Context())
		if err != nil {
			setErrorReponse(w, http.StatusInternalServerError, err)
			return
		}

		setJSONResponse(w, http.StatusOK, events)
	}
}

func getEvent(cmps *components) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			setErrorReponse(w, http.StatusBadRequest, err)
			return
		}

		event, err := cmps.App.GetEvent(r.Context(), id)
		if err != nil {
			setErrorReponse(w, getErrorStatus(err), err)
			return
		}

		setJSONResponse(w, http.StatusOK, event)
	}
}

func createEvent(cmps *components) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := getPayload[models.Event](r, cmps.Validator)
		if err != nil {
			setErrorReponse(w, http.StatusBadRequest, err)
			return
		}

		createdEvent, err := cmps.App.CreateEvent(r.Context(), request)
		if err != nil {
			setErrorReponse(w, getErrorStatus(err), err)
			return
		}

		setJSONResponse(w, http.StatusCreated, createdEvent)
	}
}

func updateEvent(cmps *components) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			setErrorReponse(w, http.StatusBadRequest, err)
			return
		}

		request, err := getPayload[models.Event](r, cmps.Validator)
		if err != nil {
			setErrorReponse(w, http.StatusBadRequest, err)
			return
		}

		updatedEvent, err := cmps.App.UpdateEvent(r.Context(), id, request)
		if err != nil {
			setErrorReponse(w, getErrorStatus(err), err)
			return
		}

		setJSONResponse(w, http.StatusOK, updatedEvent)
	}
}

func deleteEvent(cmps *components) handlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			setErrorReponse(w, http.StatusBadRequest, err)
			return
		}

		if err := cmps.App.DeleteEvent(r.Context(), id); err != nil {
			setErrorReponse(w, getErrorStatus(err), err)
			return
		}

		setEmptyResponse(w, http.StatusNoContent)
	}
}
