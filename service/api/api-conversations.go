package api

import (
	"encoding/json"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	pageSize, errPageSize := validators.SanitizePageOffset(r.URL.Query().Get("size"))
	pageNumber, errPageNumber := validators.SanitizePageNumber(r.URL.Query().Get("page"))

	if errPageSize != nil || errPageNumber != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if ok, _ := validators.UsernameIsValid(username); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUser := r.Context().Value("authedUser").(*database.User)
	if authedUser.Username != username {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	conversations, err := authedUser.GetConversations(pageSize, pageNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(conversations)
	return
}
