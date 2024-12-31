package storage

import (
	"fmt"
	"io"
	"os"
)

type FilesystemStorage struct {
	rootDir string
}

func NewFileSystemStorage(rootDir string) (Storage, error) {
	return FilesystemStorage{
		rootDir,
	}, nil
}

func (fs FilesystemStorage) SaveFile(filename string, content io.Reader) (string, error) {
	fullPath := fs.rootDir + filename
	dst, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("ciao")
		return "", err
	}

	if _, err := io.Copy(dst, content); err != nil {
		if err := dst.Close(); err != nil {
			return "", err
		}
		if err := os.Remove(fullPath); err != nil {
			return "", err
		}
		return "", err
	}

	if err := dst.Close(); err != nil {
		return "", err
	}

	return fullPath, nil
}

func (fs FilesystemStorage) DeleteFile(filename string) error {
	fullPath := fs.rootDir + filename
	if err := os.Remove(fullPath); err != nil {
		return err
	}
	return nil
}

func (fs FilesystemStorage) GetFile(filename string) (*os.File, error) {
	fullPath := fs.rootDir + filename
	return os.Open(fullPath)
}

func (fs FilesystemStorage) GetFilePath(filename string) string {
	fullPath := fs.rootDir + filename
	return fullPath
}

func (fs FilesystemStorage) GetRoot() string {
	return fs.rootDir
}
