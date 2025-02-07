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

type MessageInfoRouter struct {
	router
	Controller controllers.MessageInfoController
}

func NewMessageInfoRouter(routeFactory routes.RouteFactory, controller controllers.MessageInfoController) ControllerRouter {
	result := &MessageInfoRouter{
		newBaseRouter(routeFactory),
		controller,
	}
	result.initializeRoutes()
	return result
}

func (router *MessageInfoRouter) initializeRoutes() {
	router.routes = map[string]routes.Route{
		"getMessageComments": router.routeFactory.New(
			"/users/:userUUID/conversations/:conversationId/messages/:messageId/comments",
			http.MethodGet,
			router.GetMessageComments,
			true,
		),
		"commentMessage": router.routeFactory.New(
			"/users/:userUUID/conversations/:conversationId/messages/:messageId/comments",
			http.MethodPut,
			router.SetMessageComment,
			true,
		),
		"uncommentMessage": router.routeFactory.New(
			"/users/:userUUID/conversations/:conversationId/messages/:messageId/comments",
			http.MethodDelete,
			router.RemoveMessageComment,
			true,
		),
	}
}

func (router *MessageInfoRouter) GetMessageComments(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationMessageUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return apierrors.Forbidden()
	}

	comments, err := router.Controller.GetComments(urlParams.ConversationId, urlParams.MessageId, authedUserUUID)
	if err != nil {
		return err
	}

	return views.SendJson(w, comments)
}

func (router *MessageInfoRouter) SetMessageComment(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationMessageUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return apierrors.Forbidden()
	}

	requestBody := CommentRequestBody{}
	if err := parsers.ParseAndValidateRequestBody(r, &requestBody); err != nil {
		return err
	}

	comment, err := router.Controller.CommentMessage(urlParams.ConversationId, urlParams.MessageId, authedUserUUID, requestBody.Comment)
	if err != nil {
		return err
	}

	return views.SendJson(w, comment)
}

func (router *MessageInfoRouter) RemoveMessageComment(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationMessageUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return apierrors.Forbidden()
	}

	if err := router.Controller.UncommentMessage(urlParams.ConversationId, urlParams.MessageId, authedUserUUID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
