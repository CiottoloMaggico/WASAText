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

type SerializedImage struct {
	Uuid    string `json:"uuid"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	FullUrl string `json:"fullUrl"`
}

type Image struct {
	Uuid      string `json:"uuid"`
	Extension string `json:"-"`

	db *appdbimpl
}

func (i Image) Filename() string {
	return i.Uuid + i.Extension
}

func (i Image) FullUrl() string {
	// TODO: fix with configuration url
	return "http://127.0.0.1:3000/media/" + i.Filename()
}

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

func (i Image) Validate() error {
	if strings.HasPrefix(i.Uuid, "default_") {
		return nil
	}
	if _, err := uuid.Parse(i.Uuid); err != nil {
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
		image.Uuid,
		image.Extension,
		fileHeader.Size,
	); err != nil {
		if err := os.Remove("./media/images/" + image.Filename()); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to save the image in the db: %w", err)
	}

	return &image, nil
}

func (db *appdbimpl) GetImage(UUID string) (*Image, error) {
	image := Image{}

	if err := db.c.QueryRow(qGetImage, UUID).Scan(
		&image.Uuid,
		&image.Extension,
	); err != nil {
		return nil, err
	}

	image.db = db
	return &image, nil
}
