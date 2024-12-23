package controllers

import (
	"database/sql"
	"errors"
	"github.com/ciottolomaggico/wasatext/service/controller/translators"
	"github.com/ciottolomaggico/wasatext/service/dao"
	"github.com/ciottolomaggico/wasatext/service/views"
)

type SessionController interface {
	DoLogin(username string) (views.UserView, error)
}

type SessionControllerImpl struct {
	userController UserController
	dao            dao.UserDao
}

func (controller SessionControllerImpl) DoLogin(username string) (views.UserView, error) {
	user, err := controller.dao.GetUserByUsername(username)
	if errors.Is(err, sql.ErrNoRows) {
		newUser, err := controller.dao.CreateUser(username)
		if err != nil {
			return views.UserView{}, err
		}

		user, err = controller.dao.GetUser(newUser.GetUUID())
		if err != nil {
			return views.UserView{}, err
		}
	}

	return translators.UserWithImageToView(*user), nil
}
