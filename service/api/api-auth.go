package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/utils/authentication"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Parse the request body to get the username
	var requestBody struct {
		Username string `json:"username"`
	}
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
	user, err := rt.db.GetUser(requestBody.Username)
	if errors.Is(err, sql.ErrNoRows) {
		// create the user
		user, err = database.NewUser(requestBody.Username, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		user.Save(rt.db)
	}

	// Return the user data to the client
	json.NewEncoder(w).Encode(user)
}

func GetAuthenticatedUser(r *http.Request, db database.AppDatabase) (*database.User, error) {
	token := authentication.GetAuthToken(r)
	if _, err := uuid.Parse(token); err != nil {
		return nil, errors.New("Invalid token or User not authenticated")
	}

	user, err := db.GetUserByUUID(token)
	if err != nil {
		return nil, errors.New("The user does not exist")
	}
	return user, err
}
