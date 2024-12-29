package routers

import (
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	controllers "github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/ggicci/httpin"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"path/filepath"
)

type MessageRouter struct {
	Controller controllers.MessageController
}

func (router MessageRouter) ListRoutes() []routes.Route {
	return []routes.Route{
		routes.New(
			"/users/:UserUUID/conversations/:ConversationId/messages",
			http.MethodPost,
			router.SendMessage,
			true,
		),
		routes.New(
			"/users/:UserUUID/conversations/:ConversationId/messages",
			http.MethodPut,
			router.SetSeen,
			true,
		),
		routes.New(
			"/users/:UserUUID/conversations/:ConversationId/messages",
			http.MethodGet,
			router.GetConversationMessages,
			true,
		),
		routes.New(
			"/users/:UserUUID/conversations/:ConversationId/messages/:MessageId",
			http.MethodGet,
			router.GetConversationMessageDetail,
			true,
		),
		routes.New(
			"/users/:UserUUID/conversations/:ConversationId/messages/:MessageId",
			http.MethodDelete,
			router.DeleteConversationMessage,
			true,
		),
		routes.New(
			"/users/:UserUUID/conversations/:ConversationId/messages/:MessageId/comments",
			http.MethodGet,
			router.GetMessageComments,
			true,
		),
		routes.New(
			"/users/:UserUUID/conversations/:ConversationId/messages/:MessageId/comments",
			http.MethodPut,
			router.SetMessageComment,
			true,
		),
		routes.New(
			"/users/:UserUUID/conversations/:ConversationId/messages/:MessageId/comments",
			http.MethodDelete,
			router.RemoveMessageComment,
			true,
		),
		routes.New(
			"/users/:UserUUID/conversations/:ConversationId/messages/:MessageId/forward",
			http.MethodPost,
			router.ForwardMessage,
			true,
		),
	}
}

func (router MessageRouter) SendMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserConversationUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	requestBody := r.Context().Value(httpin.Input).(*NewMessageRequestBody)
	if err := validate.Struct(requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var file *io.Reader = nil
	var fileExt *string = nil

	if requestBody.Attachment != nil {
		file, err := requestBody.Attachment.Open()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		if err := validators.ImageIsValid(requestBody.Attachment.Filename(), requestBody.Attachment.Size(), file); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tmpExt := filepath.Ext(requestBody.Attachment.Filename())
		fileExt = &tmpExt
	}

	message, err := router.Controller.SendMessage(urlParams.ConversationId, authedUserUUID, requestBody.ReplyTo, requestBody.Content, fileExt, file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	views.SendJson(w, message)
}

func (router MessageRouter) SetSeen(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserConversationUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: handle pagination in SetConversationMessagesAsSeen
	_, err := ParseAndValidatePaginationParams(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	messages, err := router.Controller.SetConversationMessagesAsSeen(urlParams.ConversationId, authedUserUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, messages)
}

func (router MessageRouter) GetConversationMessages(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserConversationUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	queryParams, err := ParseAndValidatePaginationParams(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	messages, err := router.Controller.GetConversationMessages(urlParams.ConversationId, queryParams.Page, queryParams.Size, authedUserUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, messages)
}

func (router MessageRouter) GetConversationMessageDetail(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserConversationMessageUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	message, err := router.Controller.GetConversationMessage(urlParams.ConversationId, urlParams.MessageId, authedUserUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, message)
}

func (router MessageRouter) DeleteConversationMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserConversationMessageUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//TODO: add check message author is the request issuer
	if err := router.Controller.DeleteMessage(urlParams.ConversationId, urlParams.MessageId, authedUserUUID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (router MessageRouter) GetMessageComments(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserConversationMessageUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	comments, err := router.Controller.GetComments(urlParams.ConversationId, urlParams.MessageId, authedUserUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, comments)
}

func (router MessageRouter) SetMessageComment(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserConversationMessageUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	requestBody := CommentRequestBody{}
	if err := ParseAndValidateRequestBody(r, &requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comment, err := router.Controller.CommentMessage(urlParams.ConversationId, urlParams.MessageId, authedUserUUID, requestBody.Comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, comment)
}

func (router MessageRouter) RemoveMessageComment(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserConversationMessageUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := router.Controller.UncommentMessage(urlParams.ConversationId, urlParams.MessageId, authedUserUUID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (router MessageRouter) ForwardMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserConversationMessageUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	requestBody := ForwardRequestBody{}
	if err := ParseAndValidateRequestBody(r, &requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message, err := router.Controller.ForwardMessage(urlParams.ConversationId, urlParams.MessageId, authedUserUUID, requestBody.ForwardToId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, message)
}
