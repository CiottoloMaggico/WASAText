package api_errors

import (
	"net/http"
	"strconv"
)

type ApiError struct {
	Status int         `json:"status"`
	Detail interface{} `json:"detail"`
}

func (e ApiError) Error() string {
	return strconv.Itoa(e.Status)
}

func NewApiError(statusCode int, detail interface{}) ApiError {
	return ApiError{
		statusCode,
		detail,
	}
}

func InvalidJson() error {
	return NewApiError(http.StatusBadRequest, "Invalid json body")
}

func InvalidMultipartBody() error {
	return NewApiError(http.StatusBadRequest, "Invalid multipart body")
}

func AuthenticationRequired() error {
	return NewApiError(http.StatusUnauthorized, "Authentication required")
}

func InvalidUrlParameters() error {
	return NewApiError(http.StatusBadRequest, "Invalid url parameters")
}

func UnprocessableContent(errors map[string]string) error {
	return NewApiError(http.StatusUnprocessableEntity, errors)
}

func Forbidden() error {
	return NewApiError(http.StatusForbidden, "You are not allowed to access this resource")
}

func Conflict(errors map[string]string) error {
	return NewApiError(http.StatusConflict, errors)
}

func ResourceNotFound() error {
	return NewApiError(http.StatusNotFound, "Resource not found")
}
