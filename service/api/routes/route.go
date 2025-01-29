package routes

type Route interface {
	GetPath() string
	GetMethod() string
	GetHandler() Handler
	AuthenticationRequired() bool
}

type RouteImpl struct {
	Path           string
	Method         string
	Authentication bool
	Handler        Handler
}

func (r RouteImpl) GetPath() string {
	return r.Path
}

func (r RouteImpl) GetMethod() string {
	return r.Method
}

func (r RouteImpl) AuthenticationRequired() bool {
	return r.Authentication
}

func (r RouteImpl) GetHandler() Handler {
	return r.Handler
}

func New(path string, method string, handler Handler, authRequired bool) Route {
	return RouteImpl{
		path,
		method,
		authRequired,
		handler,
	}
}
