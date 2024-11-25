package authentication

import (
	"net/http"
	"strings"
)

func GetAuthToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "bearer ")
}
