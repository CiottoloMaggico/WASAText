package controllers

import (
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
)

type MessageInfoController interface {
	GetComments(conversationId int64, messageId int64, requestIssuerUUID string) ([]views.CommentWithAuthorView, error)
	CommentMessage(conversationId int64, messageId int64, authorUUID string, comment string) (views.CommentView, error)
	UncommentMessage(conversationId int64, messageId int64, authorUUID string) error
}

type MessageInfoControllerImpl struct {
	Model             models.MessageInfoModel
	MessageController MessageController
}

func (controller MessageInfoControllerImpl) GetComments(conversationId int64, messageId int64, requestIssuerUUID string) ([]views.CommentWithAuthorView, error) {
	if _, err := controller.MessageController.GetConversationMessage(conversationId, messageId, requestIssuerUUID); err != nil {
		return nil, err
	}

	commentsCount, err := controller.Model.CountMessageComments(messageId)
	if err != nil {
		return nil, err
	}

	comments := make([]models.MessageInfoWithUser, 0)
	if commentsCount > 0 {
		comments, err = controller.Model.GetMessageComments(messageId)
		if err != nil {
			return nil, translators.DBErrorToApiError(err)
		}
	}

	return translators.MessageInfoWithUserListToCommentView(comments), nil
}

func (controller MessageInfoControllerImpl) CommentMessage(conversationId int64, messageId int64, authorUUID string, comment string) (views.CommentView, error) {
	if _, err := controller.MessageController.GetConversationMessage(conversationId, messageId, authorUUID); err != nil {
		return views.CommentView{}, err
	}

	res, err := controller.Model.SetComment(authorUUID, messageId, comment)
	if err != nil {
		return views.CommentView{}, err
	}
	return translators.MessageInfoToCommentView(*res), nil
}

func (controller MessageInfoControllerImpl) UncommentMessage(conversationId int64, messageId int64, authorUUID string) error {
	if _, err := controller.MessageController.GetConversationMessage(conversationId, messageId, authorUUID); err != nil {
		return err
	}

	if err := controller.Model.RemoveComment(authorUUID, messageId); err != nil {
		return translators.DBErrorToApiError(err)
	}

	return nil
}
