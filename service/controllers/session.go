package controllers

import (
	"errors"
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
)

type SessionController interface {
	DoLogin(username string) (views.UserView, error)
}

type SessionControllerImpl struct {
	Model models.UserModel
}

func (controller SessionControllerImpl) DoLogin(username string) (views.UserView, error) {
	user, err := controller.Model.GetUserWithImageByUsername(username)

	if errors.Is(err, database.NoResult) {
		newUser, err := controller.Model.CreateUser(username)
		if err != nil {
			return views.UserView{}, err
		}

		user, err = controller.Model.GetUserWithImage(newUser.Uuid)

		if err != nil {
			return views.UserView{}, err
		}
	}

	return translators.UserWithImageToView(*user), nil
}
