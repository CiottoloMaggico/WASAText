package api

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) ResponseHeaderMiddleware(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r, ps)
	})
}

func (rt *_router) AuthenticationMiddleware(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authedUser, err := rt.db.GetAuthenticatedUser(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "authedUser", authedUser)
		newReq := r.WithContext(ctx)
		next(w, newReq, ps)
	})
}
