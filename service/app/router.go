package app

import (
	"github.com/ciottolomaggico/wasatext/service/routers"
)

func (app *App) startRouters() map[string]interface{} {
	return map[string]interface{}{
		"user":             app.createUserRouter(),
		"session":          app.createSessionRouter(),
		"conversation":     app.createConversationRouter(),
		"userConversation": app.createUserConversationRouter(),
		"message":          app.createMessageRouter(),
		"messageInfo":      app.createMessageInfoRouter(),
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

func (app *App) GetUserRouter() routers.UserRouter {
	return app.routers["user"].(routers.UserRouter)
}

func (app *App) GetSessionRouter() routers.SessionRouter {
	router := app.routers["session"].(routers.SessionRouter)
	return router
}

func (app *App) GetConversationRouter() routers.ConversationRouter {
	return app.routers["conversation"].(routers.ConversationRouter)
}

func (app *App) GetUserConversationRouter() routers.UserConversationRouter {
	return app.routers["userConversation"].(routers.UserConversationRouter)
}

func (app *App) GetMessageRouter() routers.MessageRouter {
	return app.routers["message"].(routers.MessageRouter)
}

func (app *App) GetMessageInfoRouter() routers.MessageInfoRouter {
	return app.routers["messageInfo"].(routers.MessageInfoRouter)
}
