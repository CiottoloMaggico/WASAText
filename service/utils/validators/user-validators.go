package validators

import (
	"errors"
	"unicode/utf8"
)

func UsernameIsValid(username string) (bool, error) {
	if utf8.RuneCountInString(username) < 3 || utf8.RuneCountInString(username) > 16 {
		return false, errors.New("username length must be between 3 and 16")
	}
	return true, nil
}
