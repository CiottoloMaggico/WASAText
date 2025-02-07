package routers

import (
	"github.com/ciottolomaggico/wasatext/service/api/parsers"
	"github.com/ciottolomaggico/wasatext/service/api/requests"
	"github.com/ciottolomaggico/wasatext/service/app/routes"
	controllers "github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type SessionRouter struct {
	Router
	Controller controllers.SessionController
}

func NewSessionRouter(routeFactory routes.RouteFactory, controller controllers.SessionController) ControllerRouter {
	result := &SessionRouter{
		NewBaseRouter(routeFactory),
		controller,
	}
	result.initializeRoutes()
	return *result
}

func (router *SessionRouter) initializeRoutes() {
	router.routes = map[string]routes.Route{
		"doLogin": router.routeFactory.New(
			"/session",
			http.MethodPost,
			router.DoLogin,
			false,
		),
	}
}

func (router SessionRouter) DoLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
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
