package validators

import (
	"fmt"
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/google/uuid"
)

func UsersUUIDValidator(users []string) (bool, error) {
	for _, userUUID := range users {
		if _, err := uuid.Parse(userUUID); err != nil {
			return false, api_errors.UnprocessableContent(map[string]string{"users": fmt.Sprintf("%s isn't a valid v4 uuid", userUUID)})
		}
	}
	return true, nil
}
