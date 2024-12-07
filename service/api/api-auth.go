package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the request body to get the username
	requestBody := UsernameRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if ok, _ := validators.UsernameIsValid(requestBody.Username); !ok {
		// TODO: add correct error explanation
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse the user from the database, if it not exist then create a new user with the provided username
	user, err := rt.db.GetUserByUsername(requestBody.Username)
	if errors.Is(err, sql.ErrNoRows) {
		// create the user
		user, err = rt.db.NewUser(requestBody.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Return the user data to the client
	json.NewEncoder(w).Encode(user)
}
