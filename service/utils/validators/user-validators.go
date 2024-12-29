package validators

import (
	"fmt"
	"github.com/google/uuid"
)

func UsersUUIDValidator(users []string) (bool, error) {
	for _, userUUID := range users {
		if _, err := uuid.Parse(userUUID); err != nil {
			return false, fmt.Errorf("invalid user UUID: %s", userUUID)
		}
	}
	return true, nil
}
