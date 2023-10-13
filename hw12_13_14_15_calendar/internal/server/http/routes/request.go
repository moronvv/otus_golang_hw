package internalhttproutes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, fmt.Errorf("invalid ID value in url path")
	}

	return int64(id), nil
}

func getPayload[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	var v T
	if err := decoder.Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}
