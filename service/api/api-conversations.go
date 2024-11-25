package api

//
//import (
//	"github.com/ciottolomaggico/wasatext/service/database"
//	"github.com/ciottolomaggico/wasatext/service/utils/validators"
//	"github.com/julienschmidt/httprouter"
//	"net/http"
//)
//
//type ConversationsPaginator struct {
//	Pagination    Page                    `json:"pagination"`
//	Conversations []database.Conversation `json:"users"`
//}
//
//func (rt *_router) getMyConversations(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
//	w.Header().Set("Content-Type", "application/json")
//
//	// Check authentication
//	authedUser, err := GetAuthenticatedUser(req, rt.db)
//	if err != nil {
//		w.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//
//	// Retrieve url parameters (user_id)
//	username := ps.ByName("username")
//	if ok, _ := validators.UsernameIsValid(username); !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// Check if the authenticated user match with the provided user_id
//	if authedUser.Username != username {
//		w.WriteHeader(http.StatusForbidden)
//		return
//	}
//
//	// Retrieve the conversations ordered by latestMessage desc
//	convs, err := rt.db.GetConversations(authedUser.Username)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//
//	// Setup pagination
//
//	// return all the conversations
//}
//
//func (rt *_router) setDelivered(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {}
//
//func (rt *_router) getConversation(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {}
