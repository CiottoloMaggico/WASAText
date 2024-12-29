package controllers

import (
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
	"io"
)

type UserController interface {
	CreateUser(username string) (views.UserView, error)
	SetMyUsername(userUUID string, newUsername string) (views.UserView, error)
	SetMyPhoto(userUUID string, photoExtension string, photoFile io.Reader) (views.UserView, error)
	GetUser(userUUID string) (views.UserView, error)
	GetUserByUsername(username string) (views.UserView, error)
	GetUsers(page uint, pageSize uint) ([]views.UserView, error)
}

type UserControllerImpl struct {
	ImageController ImageController
	Model           models.UserModel
}

func (controller UserControllerImpl) CreateUser(username string) (views.UserView, error) {
	if _, err := controller.Model.CreateUser(username); err != nil {
		return views.UserView{}, err
	}

	return controller.GetUser(username)
}

func (controller UserControllerImpl) SetMyUsername(userUUID string, newUsername string) (views.UserView, error) {
	_, err := controller.Model.UpdateUsername(userUUID, newUsername)
	if err != nil {
		return views.UserView{}, err
	}

	user, err := controller.Model.GetUserWithImage(userUUID)
	if err != nil {
		return views.UserView{}, err
	}
	return translators.UserWithImageToView(*user), nil
}

func (controller UserControllerImpl) SetMyPhoto(userUUID string, photoExtension string, photoFile io.Reader) (views.UserView, error) {
	image, err := controller.ImageController.CreateImage(photoExtension, photoFile)
	if err != nil {
		return views.UserView{}, err
	}

	_, err = controller.Model.UpdateProfilePic(userUUID, image.Uuid)
	if err != nil {
		return views.UserView{}, err
	}

	user, err := controller.Model.GetUserWithImage(userUUID)
	if err != nil {
		return views.UserView{}, err
	}
	return translators.UserWithImageToView(*user), nil
}

func (controller UserControllerImpl) GetUser(userUUID string) (views.UserView, error) {
	user, err := controller.Model.GetUserWithImage(userUUID)
	if err != nil {
		return views.UserView{}, err
	}

	return translators.UserWithImageToView(*user), nil
}

func (controller UserControllerImpl) GetUserByUsername(username string) (views.UserView, error) {
	user, err := controller.Model.GetUserWithImageByUsername(username)
	if err != nil {
		return views.UserView{}, err
	}
	return translators.UserWithImageToView(*user), nil
}

func (controller UserControllerImpl) GetUsers(page uint, pageSize uint) ([]views.UserView, error) {
	users, err := controller.Model.GetUsersWithImage(page, pageSize)
	if err != nil {
		return []views.UserView{}, err
	}
	return translators.UserWithImageListToView(users), nil
}
