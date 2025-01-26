package middlewares

import (
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Model models.UserModel
}

func (m AuthMiddleware) Wrap(next routes.Handler) routes.Handler {
	return routes.Handler(func(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
		authHeader := strings.ToLower(r.Header.Get("Authorization"))
		token := strings.TrimPrefix(authHeader, "bearer ")

		if token == authHeader {
			return api_errors.AuthenticationRequired()
		}

		authedUser, err := m.Model.GetUserWithImage(token)
		if err != nil {
			return err
		}
		//if errors.Is(err, database.NoResult) {
		//	return api_errors.AuthenticationRequired()
		//}
		context.IssuerUUID = &authedUser.Uuid
		return next(w, r, params, context)
	})
}
