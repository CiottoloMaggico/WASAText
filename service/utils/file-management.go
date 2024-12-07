package utils

import (
	"io"
	"os"
)

func SaveFile(path string, content io.Reader) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, content); err != nil {
		if err := os.Remove(path); err != nil {
			return err
		}
		return err
	}
	if err := dst.Close(); err != nil {
		return err
	}
	return nil
}
