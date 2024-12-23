package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ciottolomaggico/wasatext/service/utils"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
	"github.com/google/uuid"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

var ImageNotExist = errors.New("The image doesn't exist")

//type SerializedImage struct {
//	Uuid    string `json:"uuid"`
//	Width   int    `json:"width"`
//	Height  int    `json:"height"`
//	FullUrl string `json:"fullUrl"`
//}

func (i *Image) MarshalJSON() ([]byte, error) {
	file, err := os.Open("./media/images/" + i.Filename())
	if err != nil {
		return nil, err
	}
	defer file.Close()
	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	return json.Marshal(&SerializedImage{
		i.Uuid,
		cfg.Width,
		cfg.Height,
		i.FullUrl(),
	})
}

type Image struct {
	uuid      string
	extension string

	db *appdbimpl
}

func (i Image) GetUUID() string {
	return i.uuid
}

func (i Image) GetExtension() string {
	return i.extension
}

func (i Image) Filename() string {
	return i.uuid + i.extension
}

func (i Image) FullUrl() string {
	// TODO: fix with configuration url
	return "http://127.0.0.1:3000/media/" + i.Filename()
}

func (i Image) Validate() error {
	if strings.HasPrefix(i.uuid, "default_") {
		return nil
	}
	if _, err := uuid.Parse(i.uuid); err != nil {
		return err
	}
	if _, err := os.Stat("./media/images/" + i.Filename()); errors.Is(err, os.ErrNotExist) {
		return ImageNotExist
	}

	return nil
}

// TODO: better error handling
func (db *appdbimpl) NewImage(fileHeader multipart.FileHeader, file multipart.File) (*Image, error) {
	// For each provided arguments run the corresponding validator
	if err := validators.ImageIsValid(fileHeader, file); err != nil {
		return nil, err
	}
	// If all the arguments are valid then set the "private" object fields (e.g. primary key)
	rawUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate a new uuid: %w", err)
	}
	image := Image{rawUUID.String(), filepath.Ext(fileHeader.Filename), db}

	// Save the image into the filesystem
	if err := utils.SaveFile("./media/images/"+image.Filename(), file); err != nil {
		return nil, err
	}

	// Create the image in the database
	if _, err := db.c.Exec(qCreateImage,
		image.uuid,
		image.extension,
	); err != nil {
		if err := os.Remove("./media/images/" + image.Filename()); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to save the image in the db: %w", err)
	}

	return &image, nil
}

func (db *appdbimpl) GetImage(UUID string) (*Image, error) {
	image := Image{db: db}

	if err := db.c.QueryRow(qGetImage, UUID).Scan(
		&image.uuid,
		&image.extension,
	); err != nil {
		return nil, err
	}

	return &image, nil
}
