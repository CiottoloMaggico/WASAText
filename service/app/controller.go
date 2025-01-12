package app

import (
	controllers "github.com/ciottolomaggico/wasatext/service/controllers"
	"github.com/ciottolomaggico/wasatext/service/controllers/filters"
)

func (app *App) createConversationController() controllers.ConversationController {
	return controllers.ConversationControllerImpl{
		app.createConversationModel(),
		app.createImageController(),
		app.createUserConversationController(),
	}
}

func (app *App) createUserConversationController() controllers.UserConversationController {
	filter, err := filters.NewConversationFilter()
	if err != nil {
		panic(err)
	}
	return controllers.UserConversationControllerImpl{
		app.createUserConversationModel(),
		app.createConversationModel(),
		filter,
	}
}

func (app *App) createMessageController() controllers.MessageController {
	filter, err := filters.NewMessageFilter()
	if err != nil {
		panic(err)
	}
	return controllers.MessageControllerImpl{
		app.createMessageModel(),
		app.createConversationController(),
		app.createUserConversationController(),
		app.createImageController(),
		filter,
	}
}

func (app *App) createMessageInfoController() controllers.MessageInfoController {
	return controllers.MessageInfoControllerImpl{
		app.createMessageInfoModel(),
		app.createMessageController(),
	}
}

func (app *App) createImageController() controllers.ImageController {
	return controllers.ImageControllerImpl{
		app.createImageModel(),
	}

}

func (app *App) createUserController() controllers.UserController {
	filter, err := filters.NewUserFilter()
	if err != nil {
		panic(err)
	}
	return controllers.UserControllerImpl{
		app.createImageController(),
		app.createUserModel(),
		filter,
	}
}

func (app *App) createSessionController() controllers.SessionController {
	return controllers.SessionControllerImpl{
		app.createUserModel(),
		app.createUserController(),
	}
}
