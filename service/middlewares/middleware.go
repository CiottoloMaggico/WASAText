package middlewares

import (
	"github.com/ciottolomaggico/wasatext/service/api/routes"
)

type Middleware interface {
	Wrap(next routes.Handler) routes.Handler
}
