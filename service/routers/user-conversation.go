package routers

import (
	apierrors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/parsers"
	"github.com/ciottolomaggico/wasatext/service/api/requests"
	"github.com/ciottolomaggico/wasatext/service/app/routes"
	"github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type UserConversationRouter struct {
	router
	Controller controllers.UserConversationController
}

func NewUserConversationRouter(routeFactory routes.RouteFactory, controller controllers.UserConversationController) ControllerRouter {
	result := &UserConversationRouter{
		newBaseRouter(routeFactory),
		controller,
	}
	result.initializeRoutes()
	return result
}

func (router *UserConversationRouter) initializeRoutes() {
	router.routes = map[string]routes.Route{
		"getMyConversations": router.routeFactory.New(
			"/users/:userUUID/conversations",
			http.MethodGet,
			router.GetMyConversations,
			true,
		),
		"getConversation": router.routeFactory.New(
			"/users/:userUUID/conversations/:conversationId",
			http.MethodGet,
			router.GetConversation,
			true,
		),
	}
}

func (router *UserConversationRouter) GetMyConversations(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	paginationParams, err := parsers.ParseAndValidatePaginationParams(r.URL)
	if err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return apierrors.Forbidden()
	}

	conversations, err := router.Controller.GetUserConversations(authedUserUUID, paginationParams)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversations)
}

func (router *UserConversationRouter) GetConversation(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return apierrors.Forbidden()
	}

	conversation, err := router.Controller.GetUserConversation(authedUserUUID, urlParams.ConversationId)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}
