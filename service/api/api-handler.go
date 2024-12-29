package api

import (
	"github.com/ciottolomaggico/wasatext/service/routers"
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (r *_router) Handler(routers []routers.ControllerRouter) http.Handler {
	for _, router := range routers {
		for _, route := range router.ListRoutes() {
			r.Register(route)
		}
	}

	return r.router
}
