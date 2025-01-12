package app

import (
	"github.com/ciottolomaggico/wasatext/service/routers"
)

func (app *App) StartRouters() []routers.ControllerRouter {
	return []routers.ControllerRouter{
		app.createUserRouter(),
		app.createConversationRouter(),
		app.createMessageRouter(),
		app.createMessageInfoRouter(),
		app.createSessionRouter(),
		app.createUserConversationRouter(),
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
	}
}

func (app *App) createUserConversationRouter() routers.UserConversationRouter {
	return routers.UserConversationRouter{
		app.createUserConversationController(),
	}
}

func (app *App) createMessageRouter() routers.MessageRouter {
	return routers.MessageRouter{
		app.createMessageController(),
	}
}

func (app *App) createMessageInfoRouter() routers.MessageInfoRouter {
	return routers.MessageInfoRouter{
		app.createMessageInfoController(),
	}
}

func (app *App) createSessionRouter() routers.SessionRouter {
	return routers.SessionRouter{
		app.createSessionController(),
	}
}
