package views

import (
	"encoding/json"
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"net/http"
)

func SendJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}

func ThrowError(w http.ResponseWriter, err api_errors.ApiError) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(err.Status)

	if err := json.NewEncoder(w).Encode(err); err != nil {
		return err
	}

	return nil
}
