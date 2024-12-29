package app

import (
	"github.com/ciottolomaggico/wasatext/service/routers"
)

func (app *App) StartRouters() []routers.ControllerRouter {
	return []routers.ControllerRouter{
		app.createUserRouter(),
		app.createConversationRouter(),
		app.createMessageRouter(),
		app.createSessionRouter(),
	}
}

func (app *App) createUserRouter() routers.UserRouter {
	return routers.UserRouter{
		app.createUserController(),
	}
}

func (app *App) createConversationRouter() routers.ConversationRouter {
	return routers.ConversationRouter{
		app.createConversationController(),
		app.createMessageController(),
	}
}

func (app *App) createMessageRouter() routers.MessageRouter {
	return routers.MessageRouter{
		app.createMessageController(),
	}
}

func (app *App) createSessionRouter() routers.SessionRouter {
	return routers.SessionRouter{
		app.createSessionController(),
	}
}
