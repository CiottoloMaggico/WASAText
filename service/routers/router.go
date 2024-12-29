package routers

import "github.com/ciottolomaggico/wasatext/service/api/routes"

type ControllerRouter interface {
	ListRoutes() []routes.Route
}
