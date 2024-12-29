package app

import (
	"errors"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/storage"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type App struct {
	database database.AppDatabase
	storage  storage.Storage
	logger   *logrus.Logger
}

func New(rawDB *sqlx.DB, storageRootPath *string, logger *logrus.Logger) (App, error) {
	if logger == nil {
		panic("logger cannot be nil")
		return App{}, errors.New("logger cannot be nil")
	}
	if rawDB == nil {
		panic("database cannot be nil")
		return App{}, errors.New("database cannot be nil")
	}
	if storageRootPath == nil {
		panic("storageRootPath cannot be nil")
		return App{}, errors.New("storageRootPath cannot be nil")
	}

	appDB, err := database.New(rawDB)
	if err != nil {
		return App{}, err
	}
	appStorage, err := storage.NewFileSystemStorage(*storageRootPath)
	if err != nil {
		return App{}, err
	}

	return App{
		appDB,
		appStorage,
		logger,
	}, nil
}
