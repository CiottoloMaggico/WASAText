package routers

import (
	"github.com/ciottolomaggico/wasatext/service/app/routes"
)

type ControllerRouter interface {
	GetRoute(routeId string) routes.Route
}

type Router struct {
	routeFactory routes.RouteFactory
	routes       map[string]routes.Route
}

func (r Router) GetRoute(routeId string) routes.Route {
	if route, ok := r.routes[routeId]; ok {
		return route
	}
	return nil
}

func NewBaseRouter(routeFactory routes.RouteFactory) Router {
	return Router{routeFactory, nil}
}
