package routers

import (
	apierrors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/parsers"
	"github.com/ciottolomaggico/wasatext/service/api/requests"
	"github.com/ciottolomaggico/wasatext/service/app/routes"
	"github.com/ciottolomaggico/wasatext/service/controllers"
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
	router
	Controller controllers.ConversationController
}

func NewConversationRouter(routeFactory routes.RouteFactory, controller controllers.ConversationController) ControllerRouter {
	result := &ConversationRouter{
		newBaseRouter(routeFactory),
		controller,
	}
	result.initializeRoutes()
	return result
}

func (router *ConversationRouter) initializeRoutes() {
	router.routes = map[string]routes.Route{
		"createGroup": router.routeFactory.New(
			"/users/:userUUID/groups",
			http.MethodPost,
			router.CreateGroup,
			true,
		),
		"addToGroup": router.routeFactory.New(
			"/users/:userUUID/groups/:conversationId",
			http.MethodPut,
			router.AddToGroup,
			true,
		),
		"leaveGroup": router.routeFactory.New(
			"/users/:userUUID/groups/:conversationId",
			http.MethodDelete,
			router.LeaveGroup,
			true,
		),
		"setGroupName": router.routeFactory.New(
			"/users/:userUUID/groups/:conversationId/name",
			http.MethodPut,
			router.SetGroupName,
			true,
		),
		"setGroupPhoto": router.routeFactory.New(
			"/users/:userUUID/groups/:conversationId/photo",
			http.MethodPut,
			router.SetGroupPhoto,
			true,
		),
		"createChat": router.routeFactory.New(
			"/users/:userUUID/chats",
			http.MethodPost,
			router.CreateChat,
			true,
		),
	}
}

func (router *ConversationRouter) CreateGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return apierrors.Forbidden()
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
		defer func() {
			_ = file.Close()
		}()
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

func (router *ConversationRouter) AddToGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return apierrors.Forbidden()
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

func (router *ConversationRouter) LeaveGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return apierrors.Forbidden()
	}

	if err := router.Controller.LeaveGroup(urlParams.ConversationId, authedUserUUID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (router *ConversationRouter) SetGroupName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return apierrors.Forbidden()
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

func (router *ConversationRouter) SetGroupPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserConversationUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return apierrors.Forbidden()
	}

	requestBody := GroupPhotoRequestBody{}
	if err := parsers.ParseAndValidateMultipartRequestBody(r, &requestBody); err != nil {
		return err
	}
	file, err := requestBody.Photo.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	conversation, err := router.Controller.ChangeGroupPhoto(urlParams.ConversationId, authedUserUUID, filepath.Ext(requestBody.Photo.Filename), file)
	if err != nil {
		return err
	}

	return views.SendJson(w, conversation)
}

func (router *ConversationRouter) CreateChat(w http.ResponseWriter, r *http.Request, params httprouter.Params, context requests.RequestContext) error {
	urlParams := UserUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return apierrors.Forbidden()
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
