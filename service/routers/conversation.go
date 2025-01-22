package routers

import (
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/parsers"
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	controllers "github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/julienschmidt/httprouter"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type NewChatRequestBody struct {
	Recipient string `json:"recipient" validate:"required,uuid4"`
}
type NewGroupRequestBody struct {
	Name  string                `form:"name" validate:"required,min=3,max=16"`
	Photo *multipart.FileHeader `form:"image" validate:"omitempty,image"`
}

type ConversationNameRequestBody struct {
	Name string `json:"name" validate:"required,min=3,max=16"`
}

type GroupPhotoRequestBody struct {
	Photo *multipart.FileHeader `form:"image" validate:"required,image"`
}

type AddParticipantsRequestBody struct {
	Participants []string `json:"participants" validate:"required,unique,min=1,max=200,dive,uuid4"`
}

type ConversationRouter struct {
	Controller controllers.ConversationController
}

func (router ConversationRouter) ListRoutes() []routes.Route {
	return []routes.Route{
		routes.New(
			"/users/:userUUID/groups",
			http.MethodPost,
			router.createGroup,
			true,
		),
		routes.New(
			"/users/:userUUID/groups/:conversationId",
			http.MethodPut,
			router.addToGroup,
			true,
		),
		routes.New(
			"/users/:userUUID/groups/:conversationId",
			http.MethodDelete,
			router.leaveGroup,
			true,
		),
		routes.New(
			"/users/:userUUID/groups/:conversationId/name",
			http.MethodPut,
			router.setGroupName,
			true,
		),
		routes.New(
			"/users/:userUUID/groups/:conversationId/photo",
			http.MethodPut,
			router.setGroupPhoto,
			true,
		),
		routes.New(
			"/users/:userUUID/chats",
			http.MethodPost,
			router.createChat,
			true,
		),
	}
}

func (router ConversationRouter) createGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
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
		fileReader = file
		tmpExt := filepath.Ext(requestBody.Photo.Filename)
		fileExt = &tmpExt
	}

	conversation, err := router.Controller.CreateGroup(requestBody.Name, authedUserUUID, fileExt, fileReader)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}

func (router ConversationRouter) addToGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
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

	conversation, err := router.Controller.AddToGroup(urlParams.ConversationId, authedUserUUID, requestBody.Participants)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}

func (router ConversationRouter) leaveGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	if err := router.Controller.LeaveGroup(urlParams.ConversationId, authedUserUUID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (router ConversationRouter) setGroupName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
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

	conversation, err := router.Controller.ChangeGroupName(urlParams.ConversationId, authedUserUUID, requestBody.Name)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}

func (router ConversationRouter) setGroupPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
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

	conversation, err := router.Controller.ChangeGroupPhoto(urlParams.ConversationId, authedUserUUID, filepath.Ext(requestBody.Photo.Filename), file)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}

func (router ConversationRouter) createChat(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
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

	conversation, err := router.Controller.CreateChat(authedUserUUID, requestBody.Recipient)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}
