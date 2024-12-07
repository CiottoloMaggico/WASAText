package api

//
//import (
//	"encoding/json"
//	"github.com/ciottolomaggico/wasatext/service/database"
//	"github.com/ciottolomaggico/wasatext/service/utils/validators"
//	"github.com/julienschmidt/httprouter"
//	"net/http"
//	"strconv"
//)
//
//func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	// Url parameters validation
//	username := ps.ByName("username")
//	if ok, _ := validators.UsernameIsValid(username); !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	authedUser := r.Context().Value("authedUser").(*database.User)
//	if authedUser.Username != username {
//		w.WriteHeader(http.StatusForbidden)
//		return
//	}
//
//	// Body values validation
//	name := r.FormValue("name")
//	if ok, _ := validators.GroupNameIsValid(name); !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	file, fileH, err := r.FormFile("image")
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	image, err := rt.db.NewImage(*fileH, file, *authedUser)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	participants := r.FormValue("participants")
//	participantsList := make([]string, 0)
//	if err := json.Unmarshal([]byte(participants), &participantsList); err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// Create the group
//	group, err := rt.db.NewGroup(name, *image, *authedUser, participantsList)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	json.NewEncoder(w).Encode(group)
//	return
//}
//
//func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	// Url parameters validation
//	username := ps.ByName("username")
//	if ok, _ := validators.UsernameIsValid(username); !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	authedUser := r.Context().Value("authedUser").(*database.User)
//	if authedUser.Username != username {
//		w.WriteHeader(http.StatusForbidden)
//		return
//	}
//	groupId, err := strconv.ParseInt(ps.ByName("group"), 10, 64)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	requestBody := struct {
//		Participants []string `json:"participants"`
//	}{}
//	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	conv, err := authedUser.GetConversation(groupId)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	group, ok := conv.(database.GroupConversation)
//	if !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	if err := group.AddParticipants(requestBody.Participants); err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	json.NewEncoder(w).Encode(&group)
//	return
//}
//
//func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	// Url parameters validation
//	username := ps.ByName("username")
//	if ok, _ := validators.UsernameIsValid(username); !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	authedUser := r.Context().Value("authedUser").(*database.User)
//	if authedUser.Username != username {
//		w.WriteHeader(http.StatusForbidden)
//		return
//	}
//	groupId, err := strconv.ParseInt(ps.ByName("group"), 10, 64)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	conv, err := authedUser.GetConversation(groupId)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	group, ok := conv.(database.GroupConversation)
//	if !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	group.RemoveParticipant(authedUser.Uuid)
//	w.WriteHeader(http.StatusOK)
//	return
//}
//
//func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	// retrieve url parameters
//	username := ps.ByName("username")
//	if ok, _ := validators.UsernameIsValid(username); !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	groupId, err := strconv.ParseInt(ps.ByName("group"), 10, 64)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// check if username in url param is equal to the authenticated user
//	authedUser := r.Context().Value("authedUser").(*database.User)
//	if authedUser.Username != username {
//		w.WriteHeader(http.StatusForbidden)
//		return
//	}
//
//	// retrieve body that contains new username
//	requestBody := struct {
//		Name string `json:"name"`
//	}{}
//	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	if ok, _ := validators.GroupNameIsValid(requestBody.Name); !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	conv, err := authedUser.GetConversation(groupId)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	group, ok := conv.(database.GroupConversation)
//	if !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	group.Name = requestBody.Name
//	// update the username
//	if err := group.Save(); err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	// return 200
//	json.NewEncoder(w).Encode(&group)
//	return
//}
//
//func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	// retrieve url parameters
//	username := ps.ByName("username")
//	if ok, _ := validators.UsernameIsValid(username); !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	groupId, err := strconv.ParseInt(ps.ByName("group"), 10, 64)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// check if username in url param is equal to the authenticated user
//	authedUser := r.Context().Value("authedUser").(*database.User)
//	if authedUser.Username != username {
//		w.WriteHeader(http.StatusForbidden)
//		return
//	}
//
//	conv, err := authedUser.GetConversation(groupId)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	group, ok := conv.(database.GroupConversation)
//	if !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// retrieve photo from request body
//	file, fileH, err := r.FormFile("image")
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	defer file.Close()
//
//	image, err := rt.db.NewImage(*fileH, file, *authedUser)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	group.Photo = *image
//	if err := group.Save(); err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	json.NewEncoder(w).Encode(&group)
//	return
//}
