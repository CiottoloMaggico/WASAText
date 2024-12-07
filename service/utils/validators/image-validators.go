package validators

import (
	"errors"
	"fmt"
	"io"
	"math"
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
	if fileSize < 1 || fileSize > 4*int64(math.Pow10(9)) {
		return "", errors.New("image size must be between 1 and 4096 bytes")
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
	return "", errors.New("MIME type is not allowed")

}

func ImageFilenameValidator(filename string, mimeType string) error {
	// file extension match the mime type
	validExtensions, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return fmt.Errorf("error checking extensions: %w", err)
	}
	fileExtension := filepath.Ext(filename)
	for _, extension := range validExtensions {
		if extension == fileExtension {
			return nil
		}
	}
	return errors.New("invalid file extension")
}

func ImageIsValid(header multipart.FileHeader, file multipart.File) error {
	mimeType, err := ImageContentValidator(header.Size, file)
	if err != nil {
		return fmt.Errorf("error image content validation: %w", err)
	}
	if err := ImageFilenameValidator(header.Filename, mimeType); err != nil {
		return fmt.Errorf("error image filename validation: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	return nil
}
