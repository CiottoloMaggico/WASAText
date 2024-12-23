package controllers

import (
	"github.com/ciottolomaggico/wasatext/service/controller/translators"
	"github.com/ciottolomaggico/wasatext/service/dao"
	"github.com/ciottolomaggico/wasatext/service/views"
	"io"
)

type ImageController interface {
	CreateImage(extension string, file io.Reader) (views.ImageView, error)
	DeleteImage(imageUUID string) error
	GetImage(imageUUID string) (views.ImageView, error)
}

type ImageControllerImpl struct {
	dao dao.ImageDao
}

func (controller ImageControllerImpl) CreateImage(extension string, file io.Reader) (views.ImageView, error) {
	// TODO: add validation

	image, err := controller.dao.CreateImage(extension, file)
	if err != nil {
		return views.ImageView{}, err
	}

	return translators.ImageToView(*image), nil
}

func (controller ImageControllerImpl) DeleteImage(imageUUID string) error {
	_, err := controller.dao.DeleteImage(imageUUID)
	if err != nil {
		return err
	}

	return nil
}

func (controller ImageControllerImpl) GetImage(imageUUID string) (views.ImageView, error) {
	image, err := controller.dao.GetImage(imageUUID)
	if err != nil {
		return views.ImageView{}, err
	}

	return translators.ImageToView(*image), nil
}
