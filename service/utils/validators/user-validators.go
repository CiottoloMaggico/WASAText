package validators

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"unicode/utf8"
)

type Validator interface {
	isValid() bool
}

func UsernameIsValid(username string) (bool, error) {
	if length := utf8.RuneCountInString(username); length < 3 || length > 16 {
		return false, errors.New("username length must be between 3 and 16")
	}
	return true, nil
}

func GroupNameIsValid(name string) (bool, error) {
	if length := utf8.RuneCountInString(name); length < 3 || length > 16 {
		return false, errors.New("name length must be between 3 and 16")
	}
	return true, nil
}

func UsersUUIDValidator(users []string) (bool, error) {
	for _, userUUID := range users {
		if _, err := uuid.Parse(userUUID); err != nil {
			return false, fmt.Errorf("invalid user UUID: %s", userUUID)
		}
	}
	return true, nil
}

func MessageContentValidator(message string) (bool, error) {
	if utf8.RuneCountInString(message) < 1 && utf8.RuneCountInString(message) > 4096 {
		return false, errors.New("message length must be between 1 and 4096")
	}
	return true, nil
}
