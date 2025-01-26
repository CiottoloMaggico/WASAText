package controllers

import (
	"errors"
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
)

type SessionController interface {
	DoLogin(username string) (views.UserWithImageView, error)
}

type SessionControllerImpl struct {
	Model          models.UserModel
	UserController UserController
}

func (controller SessionControllerImpl) DoLogin(username string) (views.UserWithImageView, error) {
	var errDB database.ErrDB
	user, err := controller.Model.GetUserWithImageByUsername(username)

	if errors.As(err, &errDB) && errors.Is(errDB.ErrType, database.ErrNoResult) {
		return controller.UserController.CreateUser(username)
	} else if err != nil {
		return views.UserWithImageView{}, translators.ErrDBToErrApi(err)
	}

	return translators.UserWithImageToView(*user), nil
}
