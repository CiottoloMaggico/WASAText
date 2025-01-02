package api

import (
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn routes.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = routes.RequestContext{
			ReqUUID:    reqUUID,
			IssuerUUID: nil,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Call the next handler in chain (usually, the handler function for the path)
		if err = fn(w, r, ps, ctx); err != nil {
			if err, ok := err.(api_errors.APIError); ok {
				rt.baseLogger.WithError(err).Warn("An api error occurred")
				views.SendJson(w, err)
			} else {
				rt.baseLogger.WithError(err).Error("An unexpected error occurred")
				views.SendJson(w, api_errors.InternalServerError())
			}
		}

		return
	}
}
