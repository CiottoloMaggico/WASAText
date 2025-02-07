package translators

import (
	"errors"
	apierrors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/database"
	"net/http"
)

var errDescriptionMapping = map[string]apierrors.ApiError{
	"message_author_must_be_a_participant":     apierrors.NewApiError(http.StatusUnprocessableEntity, "the author of the message must be a participant of the conversation"),
	"replied_message_within_same_conversation": apierrors.NewApiError(http.StatusUnprocessableEntity, map[string]string{"repliedMessageId": "the replied message must belongs to the same conversation"}),
	"group_participants_limit":                 apierrors.NewApiError(http.StatusConflict, "group participants limit reached, you can add at most 200 participants"),
}

func DBErrorToApiError(err error) error {
	var DBError database.DBError
	if !errors.As(err, &DBError) {
		return err
	}

	if errors.Is(DBError.ErrType, database.ErrTrigger) {
		if errApi, ok := errDescriptionMapping[err.Error()]; ok {
			return errApi
		}
		return err
	} else if errors.Is(DBError.ErrType, database.ErrUnique) {
		return apierrors.NewApiError(http.StatusConflict, "this resource already exists")
	} else if errors.Is(DBError.ErrType, database.ErrForeignKey) || errors.Is(DBError.ErrType, database.ErrCheck) {
		return apierrors.NewApiError(http.StatusUnprocessableEntity, "")
	} else if errors.Is(DBError.ErrType, database.ErrNoResult) {
		return apierrors.NewApiError(http.StatusNotFound, "this resource does not exist")
	}

	return err
}
