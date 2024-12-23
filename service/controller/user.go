package controllers

import (
	"github.com/ciottolomaggico/wasatext/service/controller/translators"
	"github.com/ciottolomaggico/wasatext/service/dao"
	"github.com/ciottolomaggico/wasatext/service/views"
	"io"
)

type UserController interface {
	CreateUser(username string) (views.UserView, error)
	SetMyUsername(userUUID string, newUsername string) (views.UserView, error)
	SetMyPhoto(userUUID string, newPhoto string) (views.UserView, error)
	GetUser(userUUID string) (views.UserView, error)
	GetUserByUsername(username string) (views.UserView, error)
	GetUsers() ([]views.UserView, error)
}

type UserControllerImpl struct {
	imageController ImageController
	dao             dao.UserDao
}

func (controller UserControllerImpl) CreateUser(username string) (views.UserView, error) {
	if _, err := controller.dao.CreateUser(username); err != nil {
		return views.UserView{}, err
	}

	return controller.GetUser(username)
}

func (controller UserControllerImpl) SetMyUsername(userUUID string, newUsername string) (views.UserView, error) {
	_, err := controller.dao.UpdateUsername(userUUID, newUsername)
	if err != nil {
		return views.UserView{}, err
	}

	user, err := controller.dao.GetUser(userUUID)
	if err != nil {
		return views.UserView{}, err
	}
	return translators.UserWithImageToView(*user), nil
}

func (controller UserControllerImpl) SetMyPhoto(userUUID string, photoExtension string, photoFile io.Reader) (views.UserView, error) {
	image, err := controller.imageController.CreateImage(photoExtension, photoFile)
	if err != nil {
		return views.UserView{}, err
	}

	_, err = controller.dao.UpdateProfilePic(userUUID, image.Uuid)
	if err != nil {
		return views.UserView{}, err
	}

	user, err := controller.dao.GetUser(userUUID)
	if err != nil {
		return views.UserView{}, err
	}
	return translators.UserWithImageToView(*user), nil
}

func (controller UserControllerImpl) GetUser(userUUID string) (views.UserView, error) {
	user, err := controller.dao.GetUser(userUUID)
	if err != nil {
		return views.UserView{}, err
	}

	return translators.UserWithImageToView(*user), nil
}

func (controller UserControllerImpl) GetUserByUsername(username string) (views.UserView, error) {
	user, err := controller.dao.GetUserByUsername(username)
	if err != nil {
		return views.UserView{}, err
	}
	return translators.UserWithImageToView(*user), nil
}

func (controller UserControllerImpl) GetUsers(page uint, pageSize uint) ([]views.UserView, error) {
	users, err := controller.dao.GetUsers(page, pageSize)
	if err != nil {
		return []views.UserView{}, err
	}
	return translators.UserWithImageListToView(users), nil
}
