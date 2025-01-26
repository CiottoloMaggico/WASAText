package translators

import (
	"errors"
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/database"
	"net/http"
)

var errDescriptionMapping = map[string]api_errors.ErrApi{
	"message_author_must_be_a_participant":     api_errors.NewApiError(http.StatusUnprocessableEntity, "the author of the message must be a participant of the conversation"),
	"replied_message_within_same_conversation": api_errors.NewApiError(http.StatusUnprocessableEntity, map[string]string{"repliedMessageId": "the replied message must belongs to the same conversation"}),
	"group_participants_limit":                 api_errors.NewApiError(http.StatusConflict, "group participants limit reached, you can add at most 200 participants"),
}

func ErrDBToErrApi(err error) error {
	var errDB database.ErrDB
	if !errors.As(err, &errDB) {
		return err
	}

	if errors.Is(errDB.ErrType, database.ErrTrigger) {
		if errApi, ok := errDescriptionMapping[err.Error()]; ok {
			return errApi
		}
		return err
	} else if errors.Is(errDB.ErrType, database.ErrUnique) {
		return api_errors.NewApiError(http.StatusConflict, "this resource already exists")
	} else if errors.Is(errDB.ErrType, database.ErrForeignKey) || errors.Is(errDB.ErrType, database.ErrCheck) {
		return api_errors.NewApiError(http.StatusUnprocessableEntity, "")
	} else if errors.Is(errDB.ErrType, database.ErrNoResult) {
		return api_errors.NewApiError(http.StatusNotFound, "this resource does not exist")
	}

	return err
}
