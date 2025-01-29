package app

import (
	"errors"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/middlewares"
	"github.com/ciottolomaggico/wasatext/service/routers"
	"github.com/ciottolomaggico/wasatext/service/storage"
	"github.com/ciottolomaggico/wasatext/service/validators"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Application interface {
	GetUserRouter() routers.UserRouter
	GetSessionRouter() routers.SessionRouter
	GetConversationRouter() routers.ConversationRouter
	GetUserConversationRouter() routers.UserConversationRouter
	GetMessageRouter() routers.MessageRouter
	GetMessageInfoRouter() routers.MessageInfoRouter

	CreateAuthMiddleware() middlewares.AuthMiddleware
}

type App struct {
	database database.AppDatabase
	storage  storage.Storage
	logger   *logrus.Logger

	routers map[string]interface{}
}

func New(rawDB *sqlx.DB, storageRootPath *string, logger *logrus.Logger) (Application, error) {
	if logger == nil {
		panic("logger cannot be nil")
		return &App{}, errors.New("logger cannot be nil")
	}
	if rawDB == nil {
		panic("database cannot be nil")
		return &App{}, errors.New("database cannot be nil")
	}
	if storageRootPath == nil {
		panic("storageRootPath cannot be nil")
		return &App{}, errors.New("storageRootPath cannot be nil")
	}

	appDB, err := database.New(rawDB)
	if err != nil {
		return &App{}, err
	}
	appStorage, err := storage.NewFileSystemStorage(*storageRootPath)
	if err != nil {
		return &App{}, err
	}
	if err := validators.NewAppValidator(); err != nil {
		return &App{}, err
	}

	app := App{
		appDB,
		appStorage,
		logger,
		nil,
	}
	app.routers = app.startRouters()

	return &app, nil
}
