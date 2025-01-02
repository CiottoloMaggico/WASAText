package routers

import (
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/parsers"
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	controllers "github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/validators"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"path/filepath"
)

type ConversationRouter struct {
	ControllerConv controllers.ConversationController
	ControllerMess controllers.MessageController
}

func (router ConversationRouter) ListRoutes() []routes.Route {
	return []routes.Route{
		routes.New(
			"/users/:userUUID/conversations",
			http.MethodGet,
			router.GetMyConversations,
			true,
		),
		routes.New(
			"/users/:userUUID/conversations",
			http.MethodPut,
			router.SetDelivered,
			true,
		),
		routes.New(
			"/users/:userUUID/conversations/:conversationId",
			http.MethodGet,
			router.GetConversation,
			true,
		),
		routes.New(
			"/users/:userUUID/groups",
			http.MethodPost,
			router.CreateGroup,
			true,
		),
		routes.New(
			"/users/:userUUID/groups/:conversationId",
			http.MethodPut,
			router.AddToGroup,
			true,
		),
		routes.New(
			"/users/:userUUID/groups/:conversationId",
			http.MethodDelete,
			router.LeaveGroup,
			true,
		),
		routes.New(
			"/users/:userUUID/groups/:conversationId/name",
			http.MethodPut,
			router.SetGroupName,
			true,
		),
		routes.New(
			"/users/:userUUID/groups/:conversationId/photo",
			http.MethodPut,
			router.SetGroupPhoto,
			true,
		),
		routes.New(
			"/users/:userUUID/chats",
			http.MethodPost,
			router.CreateChat,
			true,
		),
	}
}

func (router ConversationRouter) GetMyConversations(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
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

	conversations, err := router.ControllerConv.GetUserConversations(authedUserUUID, paginationParams)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversations)
}

func (router ConversationRouter) SetDelivered(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
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

	conversations, err := router.ControllerMess.SetAllMessageDelivered(authedUserUUID, paginationParams)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversations)
}

func (router ConversationRouter) GetConversation(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	conversation, err := router.ControllerConv.GetUserConversation(authedUserUUID, urlParams.ConversationId)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}

func (router ConversationRouter) CreateGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	requestBody := NewGroupRequestBody{}
	if err := parsers.ParseAndValidateMultipartRequestBody(r, &requestBody); err != nil {
		return err
	}

	var fileReader io.ReadSeeker
	var fileExt *string = nil
	if requestBody.Photo != nil {
		file, err := requestBody.Photo.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		if err = validators.ImageIsValid(requestBody.Photo.Filename, requestBody.Photo.Size, file); err != nil {
			return err
		}

		fileReader = file
		tmpExt := filepath.Ext(requestBody.Photo.Filename)
		fileExt = &tmpExt
	}

	conversation, err := router.ControllerConv.CreateGroup(requestBody.Name, authedUserUUID, fileExt, fileReader)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}

func (router ConversationRouter) AddToGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	requestBody := AddParticipantsRequestBody{}
	if err := parsers.ParseAndValidateRequestBody(r, &requestBody); err != nil {
		return err
	}

	if ok, err := validators.UsersUUIDValidator(requestBody.Participants); !ok {
		return err
	}

	conversation, err := router.ControllerConv.AddToGroup(urlParams.ConversationId, authedUserUUID, requestBody.Participants)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}

func (router ConversationRouter) LeaveGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	if err := router.ControllerConv.LeaveGroup(urlParams.ConversationId, authedUserUUID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (router ConversationRouter) SetGroupName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	requestBody := ConversationNameRequestBody{}
	if err := parsers.ParseAndValidateRequestBody(r, &requestBody); err != nil {
		return err
	}

	conversation, err := router.ControllerConv.ChangeGroupName(urlParams.ConversationId, authedUserUUID, requestBody.Name)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}

func (router ConversationRouter) SetGroupPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	requestBody := GroupPhotoRequestBody{}
	if err := parsers.ParseAndValidateMultipartRequestBody(r, &requestBody); err != nil {
		return err
	}
	file, err := requestBody.Photo.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	if err = validators.ImageIsValid(requestBody.Photo.Filename, requestBody.Photo.Size, file); err != nil {
		return err
	}

	conversation, err := router.ControllerConv.ChangeGroupPhoto(urlParams.ConversationId, authedUserUUID, filepath.Ext(requestBody.Photo.Filename), file)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}

func (router ConversationRouter) CreateChat(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	requestBody := NewChatRequestBody{}
	if err := parsers.ParseAndValidateRequestBody(r, &requestBody); err != nil {
		return err
	}

	conversation, err := router.ControllerConv.CreateChat(authedUserUUID, requestBody.Recipient)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}
