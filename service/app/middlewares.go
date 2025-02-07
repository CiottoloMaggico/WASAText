package app

import (
	"github.com/ciottolomaggico/wasatext/service/middlewares"
)

func (app *App) CreateAuthMiddleware() middlewares.AuthMiddleware {
	return middlewares.AuthMiddleware{
		app.createUserModel(),
	}
}
