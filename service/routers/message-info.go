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

type MessageInfoRouter struct {
	Controller controllers.MessageInfoController
}

func (router MessageInfoRouter) ListRoutes() []routes.Route {
	return []routes.Route{
		routes.New(
			"/users/:userUUID/conversations/:conversationId/messages/:messageId/comments",
			http.MethodGet,
			router.getMessageComments,
			true,
		),
		routes.New(
			"/users/:userUUID/conversations/:conversationId/messages/:messageId/comments",
			http.MethodPut,
			router.setMessageComment,
			true,
		),
		routes.New(
			"/users/:userUUID/conversations/:conversationId/messages/:messageId/comments",
			http.MethodDelete,
			router.removeMessageComment,
			true,
		),
	}
}

func (router MessageInfoRouter) getMessageComments(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserConversationMessageUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	comments, err := router.Controller.GetComments(urlParams.ConversationId, urlParams.MessageId, authedUserUUID)
	if err != nil {
		return err
	}

	return views.SendJson(w, comments)
}

func (router MessageInfoRouter) setMessageComment(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserConversationMessageUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
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

func (router MessageInfoRouter) removeMessageComment(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserConversationMessageUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	if err := router.Controller.UncommentMessage(urlParams.ConversationId, urlParams.MessageId, authedUserUUID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
