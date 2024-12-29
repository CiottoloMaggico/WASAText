package controllers

import (
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
	"io"
)

type ImageController interface {
	CreateImage(extension string, file io.Reader) (views.ImageView, error)
	DeleteImage(imageUUID string) error
	GetImage(imageUUID string) (views.ImageView, error)
}

type ImageControllerImpl struct {
	Model models.ImageModel
}

func (controller ImageControllerImpl) CreateImage(extension string, file io.Reader) (views.ImageView, error) {
	image, err := controller.Model.CreateImage(extension, file)
	if err != nil {
		return views.ImageView{}, err
	}

	return translators.ImageToView(*image), nil
}

func (controller ImageControllerImpl) DeleteImage(imageUUID string) error {
	_, err := controller.Model.DeleteImage(imageUUID)
	if err != nil {
		return err
	}

	return nil
}

func (controller ImageControllerImpl) GetImage(imageUUID string) (views.ImageView, error) {
	image, err := controller.Model.GetImage(imageUUID)
	if err != nil {
		return views.ImageView{}, err
	}

	return translators.ImageToView(*image), nil
}
