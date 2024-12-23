package dao

import (
	database "github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/storage"
	"github.com/gofrs/uuid"
	"image"
	"io"
)

type ImageDao interface {
	CreateImage(extension string, file io.Reader) (*database.Image, error)
	DeleteImage(uuid string) (*database.Image, error)
	GetImage(uuid string) (*database.Image, error)
}

type ImageDaoImpl struct {
	db      database.BaseDatabase
	storage storage.Storage
}

func (dao ImageDaoImpl) CreateImage(extension string, file io.Reader) (*database.Image, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}

	filename, width, height := newUUID.String()+extension, cfg.Width, cfg.Height
	path, err := dao.storage.SaveFile(filename, file)
	if err != nil {
		return nil, err
	}

	commit := true
	defer func() {
		if !commit {
			_ = dao.storage.DeleteFile(filename)
		}
	}()

	query := `INSERT INTO Image (uuid, extension, width, height, fullUrl)  VALUES (?, ?, ?, ?, ?) RETURNING *;`

	image, err := dao.db.QueryRowImage(query, newUUID.String(), extension, width, height, path)
	if err != nil {
		commit = false
		return nil, err
	}

	return image, nil
}

func (dao ImageDaoImpl) DeleteImage(uuid string) (*database.Image, error) {
	tx, err := dao.db.BeginTx()
	if err != nil {
		return nil, err
	}

	query := `DELETE FROM Image WHERE uuid = ? RETURNING *;`
	deletedImage, err := tx.QueryRowImage(query, uuid)
	if err != nil {
		return nil, err
	}

	if err := dao.storage.DeleteFile(deletedImage.Filename()); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return deletedImage, nil
}

func (dao ImageDaoImpl) GetImage(uuid string) (*database.Image, error) {
	query := `SELECT * FROM Image WHERE uuid = ?;`

	return dao.db.QueryRowImage(query, uuid)
}
