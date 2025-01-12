package routers

import (
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/parsers"
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	"github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type UserConversationRouter struct {
	Controller controllers.UserConversationController
}

func (router UserConversationRouter) ListRoutes() []routes.Route {
	return []routes.Route{
		routes.New(
			"/users/:userUUID/conversations",
			http.MethodGet,
			router.GetMyConversations,
			true,
		),
		routes.New(
			"/users/:userUUID/conversations/:conversationId",
			http.MethodGet,
			router.GetConversation,
			true,
		),
	}
}

func (router UserConversationRouter) GetMyConversations(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
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
		return api_errors.Forbidden()
	}

	conversations, err := router.Controller.GetUserConversations(authedUserUUID, paginationParams)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversations)
}

func (router UserConversationRouter) GetConversation(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	conversation, err := router.Controller.GetUserConversation(authedUserUUID, urlParams.ConversationId)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}
