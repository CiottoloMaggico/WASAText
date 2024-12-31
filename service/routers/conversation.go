package routers

import (
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	controllers "github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
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
			"/users/:UserUUID/conversations",
			http.MethodGet,
			router.GetMyConversations,
			true,
		),
		routes.New(
			"/users/:UserUUID/conversations",
			http.MethodPut,
			router.SetDelivered,
			true,
		),
		routes.New(
			"/users/:UserUUID/conversations/:ConversationId",
			http.MethodGet,
			router.GetConversation,
			true,
		),
		routes.New(
			"/users/:UserUUID/groups",
			http.MethodPost,
			router.CreateGroup,
			true,
		),
		routes.New(
			"/users/:UserUUID/groups/:ConversationId",
			http.MethodPut,
			router.AddToGroup,
			true,
		),
		routes.New(
			"/users/:UserUUID/groups/:ConversationId",
			http.MethodDelete,
			router.LeaveGroup,
			true,
		),
		routes.New(
			"/users/:UserUUID/groups/:ConversationId/name",
			http.MethodPut,
			router.SetGroupName,
			true,
		),
		routes.New(
			"/users/:UserUUID/groups/:ConversationId/photo",
			http.MethodPut,
			router.SetGroupPhoto,
			true,
		),
		routes.New(
			"/users/:UserUUID/chats",
			http.MethodPost,
			router.CreateChat,
			true,
		),
	}
}

func (router ConversationRouter) GetMyConversations(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	paginationParams, err := ParseAndValidatePaginationParams(r.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	conversations, err := router.ControllerConv.GetUserConversations(authedUserUUID, paginationParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, conversations)
}

func (router ConversationRouter) SetDelivered(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	paginationParams, err := ParseAndValidatePaginationParams(r.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	conversations, err := router.ControllerMess.SetAllMessageDelivered(authedUserUUID, paginationParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, conversations)
}

func (router ConversationRouter) GetConversation(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
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

	conversation, err := router.ControllerConv.GetUserConversation(authedUserUUID, urlParams.ConversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, conversation)
}

func (router ConversationRouter) CreateGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	requestBody := NewGroupRequestBody{}
	if err := ParseAndValidateMultipartRequestBody(r, &requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var fileReader io.ReadSeeker
	var fileExt *string = nil
	if requestBody.Photo != nil {
		file, err := requestBody.Photo.Open()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		if err := validators.ImageIsValid(requestBody.Photo.Filename, requestBody.Photo.Size, file); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fileReader = file
		tmpExt := filepath.Ext(requestBody.Photo.Filename)
		fileExt = &tmpExt
	}

	conversation, err := router.ControllerConv.CreateGroup(requestBody.Name, authedUserUUID, fileExt, fileReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, conversation)
}

func (router ConversationRouter) AddToGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
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

	requestBody := AddParticipantsRequestBody{}
	if err := ParseAndValidateRequestBody(r, &requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ok, _ := validators.UsersUUIDValidator(requestBody.Participants); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: add check for group participants limit in controllers
	conversation, err := router.ControllerConv.AddToGroup(urlParams.ConversationId, authedUserUUID, requestBody.Participants)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, conversation)
}

func (router ConversationRouter) LeaveGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
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

	if err := router.ControllerConv.LeaveGroup(urlParams.ConversationId, authedUserUUID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (router ConversationRouter) SetGroupName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
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

	requestBody := ConversationNameRequestBody{}
	if err := ParseAndValidateRequestBody(r, &requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conversation, err := router.ControllerConv.ChangeGroupName(urlParams.ConversationId, authedUserUUID, requestBody.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, conversation)
}

func (router ConversationRouter) SetGroupPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
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

	requestBody := GroupPhotoRequestBody{}
	if err := ParseAndValidateMultipartRequestBody(r, &requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, err := requestBody.Photo.Open()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := validators.ImageIsValid(requestBody.Photo.Filename, requestBody.Photo.Size, file); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conversation, err := router.ControllerConv.ChangeGroupPhoto(urlParams.ConversationId, authedUserUUID, filepath.Ext(requestBody.Photo.Filename), file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, conversation)
}

func (router ConversationRouter) CreateChat(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	urlParams := UserUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	requestBody := NewChatRequestBody{}
	if err := ParseAndValidateRequestBody(r, &requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conversation, err := router.ControllerConv.CreateChat(authedUserUUID, requestBody.Recipient)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, conversation)
}
