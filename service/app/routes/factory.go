package routes

import (
	"github.com/ciottolomaggico/wasatext/service/api/requests"
	"github.com/ciottolomaggico/wasatext/service/middlewares"
)

type RouteFactory struct {
	authMiddleware middlewares.AuthMiddleware
}

func NewRouteFactory(authMiddleware middlewares.AuthMiddleware) RouteFactory {
	return RouteFactory{authMiddleware: authMiddleware}
}

func (factory RouteFactory) New(path string, method string, handler requests.Handler, authenticationRequired bool) Route {
	finalHandler := handler
	if authenticationRequired {
		finalHandler = factory.authMiddleware.Wrap(finalHandler)
	}
	return RouteImpl{
		path,
		method,
		finalHandler,
	}
}
