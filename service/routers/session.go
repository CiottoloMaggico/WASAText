package routers

import (
	"github.com/ciottolomaggico/wasatext/service/api/parsers"
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
			router.doLogin,
			false,
		),
	}
}

func (router SessionRouter) doLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	requestBody := UsernameRequestBody{}
	if err := parsers.ParseAndValidateRequestBody(r, &requestBody); err != nil {
		return err
	}

	result, err := router.Controller.DoLogin(requestBody.Name)
	if err != nil {
		return err
	}

	return views.SendJson(w, result)
}
