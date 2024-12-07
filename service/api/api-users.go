package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
	"github.com/julienschmidt/httprouter"
	"math"
	"net/http"
)

type UsernameRequestBody struct {
	Username string `json:"username"`
}

type UserPaginator struct {
	Pagination Page            `json:"pagination"`
	Users      []database.User `json:"users"`
}

// TODO: remove filter from API docs
func (rt *_router) getUsers(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Get query parameters
	pageSize, errPageSize := validators.SanitizePageOffset(req.URL.Query().Get("size"))
	pageNumber, errPageNumber := validators.SanitizePageNumber(req.URL.Query().Get("page"))

	if errPageSize != nil || errPageNumber != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// retrieve the list of users based on query parameters
	users, err := rt.db.GetUsers(pageSize, pageNumber)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userCount, err := rt.db.UsersCount()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := UserPaginator{
		Pagination: *MakePage(
			pageNumber,
			int(math.Ceil(float64(userCount/pageSize)))-1,
			req.URL.RequestURI(),
		),
		Users: users,
	}
	// serialize the list and send
	json.NewEncoder(w).Encode(result)
}

func (rt *_router) getUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// retrieve url param "id"
	username := ps.ByName("username")
	if ok, _ := validators.UsernameIsValid(username); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// retrieve the user by its "username"
	user, err := rt.db.GetUserByUsername(username)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// serialize user and send
	json.NewEncoder(w).Encode(user)
	return
}

func (rt *_router) setMyUsername(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// retrieve url parameters
	username := ps.ByName("username")
	if ok, _ := validators.UsernameIsValid(username); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUser := req.Context().Value("authedUser").(*database.User)
	// check if username in url param is equal to the authenticated user
	if authedUser.Username != username {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// retrieve body that contains new username
	requestBody := UsernameRequestBody{}
	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if ok, _ := validators.UsernameIsValid(requestBody.Username); !ok {
		// TODO: add correct error explanation
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := database.UpdateUserParams{
		requestBody.Username,
		authedUser.ProfileImage,
	}
	// update the username
	if err := authedUser.Update(params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return 200
	json.NewEncoder(w).Encode(authedUser)
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// retrieve url parameters
	username := ps.ByName("username")
	if ok, _ := validators.UsernameIsValid(username); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUser := req.Context().Value("authedUser").(*database.User)
	// check if username in url param is equal to the authenticated user
	if authedUser.Username != username {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// retrieve photo from request body
	file, fileH, err := req.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	image, err := rt.db.NewImage(*fileH, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	params := database.UpdateUserParams{
		authedUser.Username,
		*image,
	}
	if err := authedUser.Update(params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(authedUser)
}
