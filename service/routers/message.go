package routers

import (
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/parsers"
	"github.com/ciottolomaggico/wasatext/service/api/requests"
	"github.com/ciottolomaggico/wasatext/service/app/routes"
	controllers "github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/julienschmidt/httprouter"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type NewMessageRequestBody struct {
	Attachment *multipart.FileHeader `form:"attachment" validate:"required_without=Content,omitempty,image"`
	Content    *string               `form:"content" validate:"required_without=Attachment,omitnil,min=0,max=4096"`
	ReplyTo    *int64                `form:"repliedMessageId" validate:"omitnil,min=0"`
}

type CommentRequestBody struct {
	Comment string `json:"comment" validate:"required,emoji"`
}

type ForwardRequestBody struct {
	ForwardToId int64 `json:"destConversationId" validate:"required,min=0"`
}

type MessageRouter struct {
	Router
	Controller controllers.MessageController
}

func NewMessageRouter(routeFactory routes.RouteFactory, controller controllers.MessageController) ControllerRouter {
	result := &MessageRouter{
		NewBaseRouter(routeFactory),
		controller,
	}
	result.initializeRoutes()
	return *result
}

func (router *MessageRouter) initializeRoutes() {
	router.routes = map[string]routes.Route{
		"setDelivered": router.routeFactory.New(
			"/users/:userUUID/conversations",
			http.MethodPut,
			router.SetDelivered,
			true,
		),
		"sendMessage": router.routeFactory.New(
			"/users/:userUUID/conversations/:conversationId/messages",
			http.MethodPost,
			router.SendMessage,
			true,
		),
		"setSeen": router.routeFactory.New(
			"/users/:userUUID/conversations/:conversationId/messages",
			http.MethodPut,
			router.SetSeen,
			true,
		),
		"getConversationMessages": router.routeFactory.New(
			"/users/:userUUID/conversations/:conversationId/messages",
			http.MethodGet,
			router.GetConversationMessages,
			true,
		),
		"messageDetail": router.routeFactory.New(
			"/users/:userUUID/conversations/:conversationId/messages/:messageId",
			http.MethodGet,
			router.GetConversationMessageDetail,
			true,
		),
		"deleteMessage": router.routeFactory.New(
			"/users/:userUUID/conversations/:conversationId/messages/:messageId",
			http.MethodDelete,
			router.DeleteConversationMessage,
			true,
		),
		"forwardMessage": router.routeFactory.New(
			"/users/:userUUID/conversations/:conversationId/messages/:messageId/forward",
			http.MethodPost,
			router.ForwardMessage,
			true,
		),
	}
}

func (router MessageRouter) SendMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	requestBody := NewMessageRequestBody{}
	if err := parsers.ParseAndValidateMultipartRequestBody(r, &requestBody); err != nil {
		return err
	}

	var fileReader io.ReadSeeker
	var fileExt *string = nil
	if requestBody.Attachment != nil {
		file, err := requestBody.Attachment.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		fileReader = file
		tmpExt := filepath.Ext(requestBody.Attachment.Filename)
		fileExt = &tmpExt
	}

	message, err := router.Controller.SendMessage(urlParams.ConversationId, authedUserUUID, requestBody.ReplyTo, requestBody.Content, fileExt, fileReader)
	if err != nil {
		return err
	}

	return views.SendJson(w, message)
}

func (router MessageRouter) SetDelivered(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
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

	conversations, err := router.Controller.SetAllMessageDelivered(authedUserUUID, paginationParams)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversations)
}

func (router MessageRouter) SetSeen(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationUrlParams{}
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

	messages, err := router.Controller.SetConversationMessagesAsSeen(urlParams.ConversationId, authedUserUUID, paginationParams)
	if err != nil {
		return err
	}

	return views.SendJson(w, messages)
}

func (router MessageRouter) GetConversationMessages(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationUrlParams{}
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

	messages, err := router.Controller.GetConversationMessages(urlParams.ConversationId, authedUserUUID, paginationParams)
	if err != nil {
		return err
	}

	return views.SendJson(w, messages)
}

func (router MessageRouter) GetConversationMessageDetail(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationMessageUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	message, err := router.Controller.GetConversationMessage(urlParams.ConversationId, urlParams.MessageId, authedUserUUID)
	if err != nil {
		return err
	}

	return views.SendJson(w, message)
}

func (router MessageRouter) DeleteConversationMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationMessageUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	if err := router.Controller.DeleteMessage(urlParams.ConversationId, urlParams.MessageId, authedUserUUID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (router MessageRouter) ForwardMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationMessageUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	requestBody := ForwardRequestBody{}
	if err := parsers.ParseAndValidateRequestBody(r, &requestBody); err != nil {
		return err
	}

	message, err := router.Controller.ForwardMessage(urlParams.ConversationId, urlParams.MessageId, authedUserUUID, requestBody.ForwardToId)
	if err != nil {
		return err
	}

	return views.SendJson(w, message)
}
