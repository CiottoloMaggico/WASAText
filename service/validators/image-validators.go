package validators

import (
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

var allowedMIMETypes = []string{
	"image/jpeg", "image/png", "image/gif",
}

func ImageContentValidator(fileSize int64, image multipart.File) (string, error) {
	// Check file size
	if fileSize < 1 || fileSize > 1<<30 {
		return "", api_errors.UnprocessableContent(map[string]string{
			"image": "image file size must be between 1 byte and 1 Gb",
		})
	}
	// Check if file content is actually an image
	buffer := make([]byte, 512) // To detect content type are read at most 512 bytes from the file
	if _, err := io.ReadFull(image, buffer); err != nil {
		return "", err
	}

	MIMEType := http.DetectContentType(buffer)
	for _, allowedType := range allowedMIMETypes {
		if allowedType == MIMEType {
			return MIMEType, nil
		}
	}
	return "", api_errors.UnprocessableContent(map[string]string{
		"image": "invalid image type, it should be one of these types: jpeg, jpg, png, gif",
	})

}

func ImageFilenameValidator(filename string, mimeType string) error {
	// file extension match the mime type
	validExtensions, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return err
	}
	fileExtension := filepath.Ext(filename)
	for _, extension := range validExtensions {
		if extension == fileExtension {
			return nil
		}
	}
	return api_errors.UnprocessableContent(map[string]string{
		"image": "invalid image extension",
	})
}

func ImageIsValid(filename string, fileSize int64, file multipart.File) error {
	mimeType, err := ImageContentValidator(fileSize, file)
	if err != nil {
		return err
	}
	if err := ImageFilenameValidator(filename, mimeType); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	return nil
}
