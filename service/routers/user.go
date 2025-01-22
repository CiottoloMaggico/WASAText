package routers

import (
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/parsers"
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	controllers "github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/julienschmidt/httprouter"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type UsernameRequestBody struct {
	Name string `json:"username" validate:"required,min=3,max=16"`
}

type UserPhotoRequestBody struct {
	Photo *multipart.FileHeader `form:"photo" validate:"required,image"`
}

type UserRouter struct {
	Controller controllers.UserController
}

func (router UserRouter) ListRoutes() []routes.Route {
	return []routes.Route{
		routes.New(
			"/users",
			http.MethodGet,
			router.getUsers,
			true,
		),
		routes.New(
			"/users/:userUUID",
			http.MethodGet,
			router.getUser,
			true,
		),
		routes.New(
			"/users/:userUUID/username",
			http.MethodPut,
			router.setMyUsername,
			true,
		),
		routes.New(
			"/users/:userUUID/avatar",
			http.MethodPut,
			router.setMyPhoto,
			true,
		),
	}
}

func (router UserRouter) getUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	// Get query parameters and validate them
	paginationParams, err := parsers.ParseAndValidatePaginationParams(r.URL)
	if err != nil {
		return err
	}

	// retrieve the list of users based on query parameters
	users, err := router.Controller.GetUsers(paginationParams)
	if err != nil {
		return err
	}

	// serialize the list and send
	return views.SendJson(w, users)
}

func (router UserRouter) getUser(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	// Get url parameters and validate them
	urlParams := UserUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	// retrieve the user by its "uuid"
	user, err := router.Controller.GetUser(urlParams.UserUUID)
	if err != nil {
		return err
	}

	// serialize user and send
	return views.SendJson(w, user)
}

func (router UserRouter) setMyUsername(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	// Get url parameters and validate them
	urlParams := UserUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	// check if "uuid" in url param correspond to the authed user
	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	// retrieve body that contains new username
	requestBody := UsernameRequestBody{}
	if err := parsers.ParseAndValidateRequestBody(r, &requestBody); err != nil {
		return err
	}

	// invoke the controllers
	updatedUser, err := router.Controller.SetMyUsername(authedUserUUID, requestBody.Name)
	if err != nil {
		return err
	}

	return views.SendJson(w, updatedUser)
}

func (router UserRouter) setMyPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	// Get url parameters and validate them
	urlParams := UserUrlParams{}
	if err := parsers.ParseAndValidateUrlParams(params, &urlParams); err != nil {
		return err
	}

	// check if "uuid" in url param correspond to the authed user
	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		return api_errors.Forbidden()
	}

	requestBody := UserPhotoRequestBody{}
	if err := parsers.ParseAndValidateMultipartRequestBody(r, &requestBody); err != nil {
		return err
	}
	file, err := requestBody.Photo.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	updatedUser, err := router.Controller.SetMyPhoto(authedUserUUID, filepath.Ext(requestBody.Photo.Filename), file)
	if err != nil {
		return err
	}

	return views.SendJson(w, updatedUser)
}
