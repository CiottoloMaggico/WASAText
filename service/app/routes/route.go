package routes

import (
	"github.com/ciottolomaggico/wasatext/service/api/requests"
)

type Route interface {
	GetPath() string
	GetMethod() string
	GetHandler() requests.Handler
}

type RouteImpl struct {
	Path    string
	Method  string
	Handler requests.Handler
}

func (r RouteImpl) GetPath() string {
	return r.Path
}

func (r RouteImpl) GetMethod() string {
	return r.Method
}

func (r RouteImpl) GetHandler() requests.Handler {
	return r.Handler
}
