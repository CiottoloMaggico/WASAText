package middlewares

import (
	"github.com/ciottolomaggico/wasatext/service/api/requests"
)

type Middleware interface {
	Wrap(next requests.Handler) requests.Handler
}
