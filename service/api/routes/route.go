package routes

type Route interface {
	GetPath() string
	GetMethod() string
	GetHandler() Handler
	AuthenticationRequired() bool
}

type RouteImpl struct {
	path           string
	method         string
	authentication bool
	handler        Handler
}

func (r RouteImpl) GetPath() string {
	return r.path
}

func (r RouteImpl) GetMethod() string {
	return r.method
}

func (r RouteImpl) AuthenticationRequired() bool {
	return r.authentication
}

func (r RouteImpl) GetHandler() Handler {
	return r.handler
}

func New(path string, method string, handler Handler, authRequired bool) Route {
	return RouteImpl{
		path:           path,
		method:         method,
		authentication: authRequired,
		handler:        handler,
	}
}
