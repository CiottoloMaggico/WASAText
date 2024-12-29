package models

import (
	database "github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/storage"
	"github.com/gofrs/uuid"
	"image"
	"io"
)

type Image struct {
	Uuid       string `db:"image_uuid"`
	Extension  string `db:"image_extension"`
	Width      int    `db:"image_width"`
	Height     int    `db:"image_height"`
	FullUrl    string `db:"image_fullUrl"`
	UploadedAt string `db:"image_uploadedAt"`
}

type ImageModel interface {
	CreateImage(extension string, file io.Reader) (*Image, error)
	DeleteImage(uuid string) (*Image, error)
	GetImage(uuid string) (*Image, error)
}

type ImageModelImpl struct {
	Db      database.AppDatabase
	Storage storage.Storage
}

func (i Image) Filename() string {
	return i.Uuid + i.Extension
}

func (model ImageModelImpl) CreateImage(extension string, file io.Reader) (*Image, error) {
	query := `
		INSERT INTO Image (uuid, extension, width, height, fullUrl)
		VALUES (?, ?, ?, ?, ?)
		RETURNING uuid image_uuid, extension image_extension, width image_width, height image_height, fullUrl image_fullUrl, uploadedAt image_uploadedAt;
	`
	newUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}
	filename, width, height := newUUID.String()+extension, cfg.Width, cfg.Height
	path := model.Storage.GetFilePath(filename)

	tx, err := model.Db.StartTx()
	if err != nil {
		return nil, err
	}

	image := Image{}
	if err := tx.QueryStructRow(&image, query, newUUID.String(), extension, width, height, path); err != nil {
		return nil, err
	}

	if _, err := model.Storage.SaveFile(filename, file); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &image, nil
}

func (model ImageModelImpl) DeleteImage(uuid string) (*Image, error) {
	query := `
		DELETE FROM Image
		WHERE uuid = ?
		RETURNING uuid image_uuid, extension image_extension, width image_width, height image_height, fullUrl image_fullUrl, uploadedAt image_uploadedAt;
	`
	tx, err := model.Db.StartTx()
	if err != nil {
		return nil, err
	}

	image := Image{}
	if err := tx.QueryStructRow(&image, query, uuid); err != nil {
		return nil, err
	}

	if err := model.Storage.DeleteFile(image.Filename()); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &image, nil
}

func (model ImageModelImpl) GetImage(uuid string) (*Image, error) {
	query := `
		SELECT uuid image_uuid, extension image_extension, width image_width, height image_height, fullUrl image_fullUrl, uploadedAt image_uploadedAt
		FROM Image
		WHERE uuid = ?;
	`

	image := Image{}
	if err := model.Db.QueryStructRow(&image, query, uuid); err != nil {
		return nil, err
	}

	return &image, nil
}
