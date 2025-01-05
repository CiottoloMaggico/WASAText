package views

import (
	"encoding/json"
	"net/http"
)

func SendJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}
