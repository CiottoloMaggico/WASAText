package validators

import (
	"errors"
	"strconv"
)

func SanitizePageNumber(pageNum string) (int, error) {
	if pageNum == "" {
		return 0, nil
	}
	result, err := strconv.Atoi(pageNum)
	if err != nil {
		return 0, errors.New("Invalid page number")
	}
	if result < 0 {
		return 0, errors.New("Invalid page number")
	}
	return result, nil
}

func SanitizePageOffset(pageOff string) (int, error) {
	if pageOff == "" {
		return 20, nil
	}
	result, err := strconv.Atoi(pageOff)
	if err != nil {
		return 0, errors.New("Invalid page offset")
	}
	if result <= 0 || result >= 20 {
		result = 20
	}
	return result, nil

}
