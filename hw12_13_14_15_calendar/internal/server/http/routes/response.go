package internalhttproutes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
)

func setEmptyResponse(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func setJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func setErrorReponse(w http.ResponseWriter, code int, err error) {
	setJSONResponse(w, code, map[string]string{"error": err.Error()})
}

func getErrorStatus(err error) int {
	switch {
	case errors.Is(err, app.ErrDocumentNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
