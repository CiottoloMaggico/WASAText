package routers

import (
	"github.com/ciottolomaggico/wasatext/service/app/routes"
)

type ControllerRouter interface {
	GetRoute(routeId string) routes.Route
}

type router struct {
	routeFactory routes.RouteFactory
	routes       map[string]routes.Route
}

func (r *router) GetRoute(routeId string) routes.Route {
	if route, ok := r.routes[routeId]; ok {
		return route
	}
	return nil
}

func newBaseRouter(routeFactory routes.RouteFactory) router {
	return router{routeFactory, nil}
}
