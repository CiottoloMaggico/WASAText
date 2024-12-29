package routers

import (
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	controllers "github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type SessionRouter struct {
	Controller controllers.SessionController
}

func (router SessionRouter) ListRoutes() []routes.Route {
	return []routes.Route{
		routes.New(
			"/session",
			http.MethodPost,
			router.DoLogin,
			false,
		),
	}
}

func (router SessionRouter) DoLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	requestBody := ConversationNameRequestBody{}
	if err := ParseAndValidateRequestBody(r, &requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := router.Controller.DoLogin(requestBody.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	views.SendJson(w, result)
}
