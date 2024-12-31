package middlewares

import (
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	"github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Controller controllers.UserController
}

func (m AuthMiddleware) Wrap(next routes.Handler) routes.Handler {
	return routes.Handler(func(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
		authHeader := strings.ToLower(r.Header.Get("Authorization"))
		token := strings.TrimPrefix(authHeader, "bearer ")

		if token == authHeader {
			w.WriteHeader(400)
			return
		}

		authedUser, err := m.Controller.GetUser(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		context.IssuerUUID = &authedUser.Uuid
		next(w, r, params, context)
	})
}
