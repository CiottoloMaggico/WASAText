package api

//
//import (
//	"database/sql"
//	"encoding/json"
//	"errors"
//	"github.com/ciottolomaggico/wasatext/service/database"
//	"github.com/ciottolomaggico/wasatext/service/utils/validators"
//	"github.com/julienschmidt/httprouter"
//	"math"
//	"net/http"
//)
//
//type UserPaginator struct {
//	Pagination Page            `json:"pagination"`
//	Users      []database.User `json:"users"`
//}
//
//// TODO: remove filter from API docs
//func (rt *_router) getUsers(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
//	w.Header().Set("Content-Type", "application/json")
//
//	// check if the user is authenticated
//	_, err := GetAuthenticatedUser(req, rt.db)
//	if err != nil {
//		w.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//
//	// Get query parameters
//	pageSize, errPageSize := validators.SanitizePageOffset(req.URL.Query().Get("size"))
//	pageNumber, errPageNumber := validators.SanitizePageNumber(req.URL.Query().Get("page"))
//
//	if errPageSize != nil || errPageNumber != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// retrieve the list of users based on query parameters
//	users, err := rt.db.GetUsers(pageSize, pageNumber)
//
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	userCount, err := rt.db.GetUsersCount()
//
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//
//	result := UserPaginator{
//		Pagination: *MakePage(
//			pageNumber,
//			int(math.Ceil(float64(userCount/pageSize)))-1,
//			req.URL.RequestURI(),
//		),
//		Users: *users,
//	}
//	// serialize the list and send
//	json.NewEncoder(w).Encode(result)
//}
//
//func (rt *_router) getUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
//	w.Header().Set("Content-Type", "application/json")
//
//	// check if the user is authenticated
//	_, err := GetAuthenticatedUser(req, rt.db)
//	if err != nil {
//		w.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//
//	// retrieve url param "id"
//	username := ps.ByName("username")
//	if ok, _ := validators.UsernameIsValid(username); !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// retrieve the user by its "id"
//	user, err := rt.db.GetUser(username)
//	if errors.Is(err, sql.ErrNoRows) {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//
//	// serialize user and send
//	json.NewEncoder(w).Encode(user)
//	return
//}
//
//func (rt *_router) setMyUsername(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
//	w.Header().Set("Content-Type", "application/json")
//
//	// check if the user is authenticated
//	authedUser, err := GetAuthenticatedUser(req, rt.db)
//	if err != nil {
//		w.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//
//	// retrieve url parameters
//	username := ps.ByName("username")
//	if ok, _ := validators.UsernameIsValid(username); !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// check if username in url param is equal to the authenticated user
//	if authedUser.Username != username {
//		w.WriteHeader(http.StatusForbidden)
//		return
//	}
//
//	// retrieve body that contains new username
//	var requestBody struct {
//		Username string `json:"username"`
//	}
//	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	if ok, _ := validators.UsernameIsValid(requestBody.Username); !ok {
//		// TODO: add correct error explanation
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// update the username
//	user, err := rt.db.UpdateUsername(authedUser.Username, requestBody.Username)
//	if errors.Is(err, sql.ErrNoRows) {
//		// create the user
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//
//	// return 200
//	json.NewEncoder(w).Encode(user)
//}
//
////func (rt *_router) setMyPhoto(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {}
