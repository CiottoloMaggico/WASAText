package routers

import (
	"encoding/json"
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	controllers "github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/ggicci/httpin"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"path/filepath"
)

// TODO: pagination hint
//type UserPaginator struct {
//	Pagination Page            `json:"pagination"`
//	Users      []database.User `json:"users"`
//}
//result := UserPaginator{
//	Pagination: *MakePage(
//		pageNumber,
//		int(math.Ceil(float64(userCount/pageSize)))-1,
//		req.URL.RequestURI(),
//	),
//	Users: users,
//}

type UserRouter struct {
	Controller controllers.UserController
}

func (router UserRouter) ListRoutes() []routes.Route {
	return []routes.Route{
		routes.New(
			"/users",
			http.MethodGet,
			router.GetUsers,
			true,
		),
		routes.New(
			"/users/:UserUUID",
			http.MethodGet,
			router.GetUser,
			true,
		),
		routes.New(
			"/users/:UserUUID/username",
			http.MethodPut,
			router.SetMyUsername,
			true,
		),
		routes.New(
			"/users/:UserUUID/avatar",
			http.MethodPut,
			router.SetMyPhoto,
			true,
		),
	}
}

func (router UserRouter) GetUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	// Get query parameters and validate them
	queryParams, err := ParseAndValidatePaginationParams(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// retrieve the list of users based on query parameters
	// TODO: handle pagination
	users, err := router.Controller.GetUsers(queryParams.Page, queryParams.Size)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// serialize the list and send
	views.SendJson(w, users)
	return
}

func (router UserRouter) GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	// Get url parameters and validate them
	urlParams := UserUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// retrieve the user by its "uuid"
	user, err := router.Controller.GetUser(urlParams.UserUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// serialize user and send
	json.NewEncoder(w).Encode(user)
	return
}

func (router UserRouter) SetMyUsername(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	// Get url parameters and validate them
	urlParams := UserUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if "uuid" in url param correspond to the authed user
	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// retrieve body that contains new username
	requestBody := UsernameRequestBody{}
	if err := ParseAndValidateRequestBody(r, &requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// invoke the controllers
	updatedUser, err := router.Controller.SetMyUsername(authedUserUUID, requestBody.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, updatedUser)
}

func (router UserRouter) SetMyPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) {
	// Get url parameters and validate them
	urlParams := UserUrlParams{}
	if err := ParseAndValidateUrlParams(params, &urlParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if "uuid" in url param correspond to the authed user
	authedUserUUID := *context.IssuerUUID
	if authedUserUUID != urlParams.UserUUID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	requestBody := r.Context().Value(httpin.Input).(*UserPhotoRequestBody)
	file, err := requestBody.Photo.Open()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := validators.ImageIsValid(requestBody.Photo.Filename(), requestBody.Photo.Size(), file); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedUser, err := router.Controller.SetMyPhoto(authedUserUUID, filepath.Ext(requestBody.Photo.Filename()), file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.SendJson(w, updatedUser)
}
