package app

import (
	"github.com/ciottolomaggico/wasatext/service/api/requests"
	"github.com/ciottolomaggico/wasatext/service/app/routes"
	"github.com/ciottolomaggico/wasatext/service/routers"
)

func (app *App) startRouters() map[string]routers.ControllerRouter {
	routeFactory := routes.NewRouteFactory(app.CreateAuthMiddleware())

	return map[string]routers.ControllerRouter{
		"user":             routers.NewUserRouter(routeFactory, app.createUserController()),
		"session":          routers.NewSessionRouter(routeFactory, app.createSessionController()),
		"conversation":     routers.NewConversationRouter(routeFactory, app.createConversationController()),
		"userConversation": routers.NewUserConversationRouter(routeFactory, app.createUserConversationController()),
		"message":          routers.NewMessageRouter(routeFactory, app.createMessageController()),
		"messageInfo":      routers.NewMessageInfoRouter(routeFactory, app.createMessageInfoController()),
	}
}

func (app *App) GetEndpointHandler(endpointId string) requests.Handler {
	for _, router := range app.routers {
		if route := router.GetRoute(endpointId); route != nil {
			return route.GetHandler()
		}
	}

	app.logger.Errorf("Endpoint not found: %s", endpointId)
	return nil
}
