package app

import "github.com/ciottolomaggico/wasatext/service/models"

func (app *App) createConversationModel() models.ConversationModel {
	return models.ConversationModelImpl{
		app.database,
	}
}

func (app *App) createUserConversationModel() models.UserConversationModel {
	return models.UserConversationModelImpl{
		app.database,
	}
}

func (app *App) createMessageModel() models.MessageModel {
	return models.MessageModelImpl{
		app.database,
	}
}

func (app *App) createMessageInfoModel() models.MessageInfoModel {
	return models.MessageInfoModelImpl{
		app.database,
	}
}

func (app *App) createImageModel() models.ImageModel {
	return models.ImageModelImpl{
		app.database,
		app.storage,
	}

}

func (app *App) createUserModel() models.UserModel {
	return models.UserModelImpl{
		app.database,
	}
}
