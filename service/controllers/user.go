package controllers

import (
	"errors"
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/ciottolomaggico/wasatext/service/views/pagination"
	"io"
	"net/http"
)

type UserController interface {
	CreateUser(username string) (views.UserWithImageView, error)
	SetMyUsername(userUUID string, newUsername string) (views.UserWithImageView, error)
	SetMyPhoto(userUUID string, photoExtension string, photoFile io.ReadSeeker) (views.UserWithImageView, error)
	GetUser(userUUID string) (views.UserWithImageView, error)
	GetUserByUsername(username string) (views.UserWithImageView, error)
	GetUsers(paginationPs pagination.PaginationParams) (pagination.PaginatedView, error)
}

type UserControllerImpl struct {
	ImageController ImageController
	Model           models.UserModel
}

func (controller UserControllerImpl) CreateUser(username string) (views.UserWithImageView, error) {
	user, err := controller.Model.CreateUser(username)
	if errors.Is(err, database.UniqueConstraint) {
		return views.UserWithImageView{}, api_errors.NewApiError(http.StatusConflict, "an user with this username already exists")
	} else if err != nil {
		return views.UserWithImageView{}, err
	}

	return controller.GetUser(user.Uuid)
}

func (controller UserControllerImpl) SetMyUsername(userUUID string, newUsername string) (views.UserWithImageView, error) {
	if _, err := controller.Model.UpdateUsername(userUUID, newUsername); errors.Is(err, database.UniqueConstraint) {
		return views.UserWithImageView{}, api_errors.NewApiError(http.StatusConflict, "an user with this username already exists")
	} else if err != nil {
		return views.UserWithImageView{}, err
	}

	return controller.GetUser(userUUID)
}

func (controller UserControllerImpl) SetMyPhoto(userUUID string, photoExtension string, photoFile io.ReadSeeker) (views.UserWithImageView, error) {
	image, err := controller.ImageController.CreateImage(photoExtension, photoFile)
	if err != nil {
		return views.UserWithImageView{}, err
	}

	if _, err = controller.Model.UpdateProfilePic(userUUID, image.Uuid); err != nil {
		return views.UserWithImageView{}, err
	}

	return controller.GetUser(userUUID)
}

func (controller UserControllerImpl) GetUser(userUUID string) (views.UserWithImageView, error) {
	user, err := controller.Model.GetUserWithImage(userUUID)
	if errors.Is(err, database.NoResult) {
		return views.UserWithImageView{}, api_errors.ResourceNotFound()
	} else if err != nil {
		return views.UserWithImageView{}, err
	}

	return translators.UserWithImageToView(*user), nil
}

func (controller UserControllerImpl) GetUserByUsername(username string) (views.UserWithImageView, error) {
	user, err := controller.Model.GetUserWithImageByUsername(username)
	if errors.Is(err, database.NoResult) {
		return views.UserWithImageView{}, api_errors.ResourceNotFound()
	} else if err != nil {
		return views.UserWithImageView{}, err
	}

	return translators.UserWithImageToView(*user), nil
}

func (controller UserControllerImpl) GetUsers(paginationPs pagination.PaginationParams) (pagination.PaginatedView, error) {
	usersCount, err := controller.Model.Count()
	if err != nil {
		return pagination.PaginatedView{}, err
	}

	if usersCount == 0 {
		return pagination.ToPaginatedView(paginationPs, usersCount, translators.UserWithImageListToView(make([]models.UserWithImage, 0)))
	}

	users, err := controller.Model.GetUsersWithImage(paginationPs.Page, paginationPs.Size)
	if err != nil {
		return pagination.PaginatedView{}, err
	}

	return pagination.ToPaginatedView(paginationPs, usersCount, translators.UserWithImageListToView(users))
}
