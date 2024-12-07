package api

//
//import (
//	"encoding/json"
//	"github.com/ciottolomaggico/wasatext/service/database"
//	"github.com/ciottolomaggico/wasatext/service/utils/validators"
//	"github.com/julienschmidt/httprouter"
//	"net/http"
//)
//
//func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	// Url parameters validation
//	username := ps.ByName("username")
//	if ok, _ := validators.UsernameIsValid(username); !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	authedUser := r.Context().Value("authedUser").(*database.User)
//	if authedUser.Username != username {
//		w.WriteHeader(http.StatusForbidden)
//		return
//	}
//
//	// Body values validation
//	requestBody := struct {
//		Recipient string `json:"recipient"`
//	}{}
//	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	recipient, err := rt.db.GetUserByUUID(requestBody.Recipient)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//
//	}
//
//	// Create the chat
//	chat, err := rt.db.NewChat(authedUser, recipient)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	json.NewEncoder(w).Encode(&chat)
//	return
//}
//
//func (rt *_router) test(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	authedUser := r.Context().Value("authedUser").(*database.User)
//	convs, _ := authedUser.GetConversation(2)
//	chat := convs.(database.ChatConversation)
//	json.NewEncoder(w).Encode(&chat)
//}
