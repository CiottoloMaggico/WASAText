package app

import (
	controllers "github.com/ciottolomaggico/wasatext/service/controllers"
)

func (app *App) createConversationController() controllers.ConversationController {
	return controllers.ConversationControllerImpl{
		app.createImageController(),
		app.createConversationModel(),
		app.createUserConversationModel(),
	}
}

func (app *App) createMessageController() controllers.MessageController {
	return controllers.MessageControllerImpl{
		app.createConversationController(),
		app.createImageController(),
		app.createMessageInfoModel(),
		app.createMessageModel(),
		app.createConversationModel(),
	}
}

func (app *App) createImageController() controllers.ImageController {
	return controllers.ImageControllerImpl{
		app.createImageModel(),
	}

}

func (app *App) createUserController() controllers.UserController {
	return controllers.UserControllerImpl{
		app.createImageController(),
		app.createUserModel(),
	}
}

func (app *App) createSessionController() controllers.SessionController {
	return controllers.SessionControllerImpl{
		app.createUserModel(),
	}
}
