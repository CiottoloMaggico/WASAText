package middlewares

import (
	"errors"
	apierrors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/requests"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Model models.UserModel
}

func (m AuthMiddleware) Wrap(next requests.Handler) requests.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
		noResult := database.ErrNoResult
		authHeader := strings.ToLower(r.Header.Get("Authorization"))
		token := strings.TrimPrefix(authHeader, "bearer ")

		if token == authHeader {
			return apierrors.AuthenticationRequired()
		}

		authedUser, err := m.Model.GetUserWithImage(token)
		if errors.As(err, &noResult) {
			return apierrors.AuthenticationRequired()
		} else if err != nil {
			return err
		}

		context.IssuerUUID = &authedUser.Uuid
		return next(w, r, params, context)
	}
}
