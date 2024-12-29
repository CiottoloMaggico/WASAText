package storage

import (
	"io"
	"os"
)

type Storage interface {
	SaveFile(filename string, content io.Reader) (string, error)
	DeleteFile(filename string) error
	GetFile(filename string) (*os.File, error)
	GetFilePath(filename string) string
	GetRoot() string
}
