package controllers

import (
	"github.com/ciottolomaggico/wasatext/service/controller/translators"
	"github.com/ciottolomaggico/wasatext/service/dao"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/views"
	"io"
)

type MessageController interface {
	SendMessage(conversationID int64, authorUUID string, replyToId *int64, content *string, attachmentExt *string, attachmentFile *io.Reader) (views.MessageView, error)
	GetConversationMessages(conversationID int64, page uint, pageSize uint, requestIssuerUUID string) ([]views.MessageView, error)
	GetConversationMessage(conversationID int64, messageId int64, requestIssuerUUID string) (views.MessageView, error)
	GetComments(conversationId int64, messageId int64, requestIssuerUUID string) ([]views.CommentView, error)
	CommentMessage(conversationId int64, messageId int64, authorUUID string, comment string) (views.CommentView, error)
	UncommentMessage(conversationId int64, messageId int64, authorUUID string) error
	ForwardMessage(conversationId int64, messageId int64, authorUUID string, forwardToConversationID int64) (views.MessageView, error)
	SetAllMessageDelivered(requestIssuerUUID string) ([]database.UserConversation, error)
	SetConversationMessagesAsSeen(conversationId int64, requestIssuer string) ([]views.MessageView, error)
}

type MessageControllerImpl struct {
	conversationController ConversationController
	userController         UserController
	imageController        ImageController
	messageInfoDao         dao.MessageInfoDao
	messageDao             dao.MessageDao
	conversationDao        dao.ConversationDao
}

func (controller MessageControllerImpl) SendMessage(conversationID int64, authorUUID string, replyToId *int64, content *string, attachmentExt *string, attachmentFile *io.Reader) (views.MessageView, error) {
	if res, err := controller.conversationDao.IsParticipant(conversationID, authorUUID); err != nil {
		return views.MessageView{}, err
	} else if !res {
		return views.MessageView{}, NotParticipant
	}

	var attachmentUUID *string = nil
	if attachmentExt != nil && attachmentFile != nil {
		image, err := controller.imageController.CreateImage(*attachmentExt, *attachmentFile)
		if err != nil {
			return views.MessageView{}, err
		}
		attachmentUUID = &image.Uuid
	}

	message, err := controller.messageDao.CreateMessage(
		conversationID, authorUUID, replyToId, content, attachmentUUID,
	)
	if err != nil {
		return views.MessageView{}, err
	}
	return controller.GetConversationMessage(conversationID, message.GetId(), authorUUID)
}

func (controller MessageControllerImpl) GetConversationMessage(conversationID int64, messageId int64, requestIssuerUUID string) (views.MessageView, error) {
	if res, err := controller.conversationDao.IsParticipant(conversationID, requestIssuerUUID); err != nil {
		return views.MessageView{}, err
	} else if !res {
		return views.MessageView{}, NotParticipant
	}

	message, err := controller.messageDao.GetConversationMessage(messageId, conversationID)
	if err != nil {
		return views.MessageView{}, err
	}
	return translators.MessageWithAuthorAndAttachmentToView(*message), err
}

func (controller MessageControllerImpl) GetConversationMessages(conversationID int64, page uint, pageSize uint, requestIssuerUUID string) ([]views.MessageView, error) {
	if res, err := controller.conversationDao.IsParticipant(conversationID, requestIssuerUUID); err != nil {
		return []views.MessageView{}, err
	} else if !res {
		return []views.MessageView{}, NotParticipant
	}

	messages, err := controller.messageDao.GetConversationMessages(conversationID, page, pageSize)
	if err != nil {
		return []views.MessageView{}, err
	}

	return translators.MessageWithAuthorAndAttachmentListToView(messages), nil
}

func (controller MessageControllerImpl) GetComments(conversationId int64, messageId int64, requestIssuerUUID string) ([]views.CommentView, error) {
	if _, err := controller.GetConversationMessage(conversationId, messageId, requestIssuerUUID); err != nil {
		return nil, err
	}

	comments, err := controller.messageInfoDao.GetMessageComments(messageId)

	if err != nil {
		return nil, err
	}

	return translators.MessageInfoListToCommentView(comments), nil
}

func (controller MessageControllerImpl) CommentMessage(conversationId int64, messageId int64, authorUUID string, comment string) (views.CommentView, error) {
	if res, err := controller.conversationDao.IsParticipant(conversationId, authorUUID); err != nil {
		return views.CommentView{}, err
	} else if !res {
		return views.CommentView{}, NotParticipant
	}

	if _, err := controller.GetConversationMessage(conversationId, messageId, authorUUID); err != nil {
		return views.CommentView{}, err
	}

	res, err := controller.messageInfoDao.SetComment(authorUUID, messageId, comment)
	if err != nil {
		return views.CommentView{}, err
	}
	return translators.MessageInfoToCommentView(*res), nil
}

func (controller MessageControllerImpl) UncommentMessage(conversationId int64, messageId int64, authorUUID string) error {
	if res, err := controller.conversationDao.IsParticipant(conversationId, authorUUID); err != nil {
		return err
	} else if !res {
		return NotParticipant
	}

	if _, err := controller.GetConversationMessage(conversationId, messageId, authorUUID); err != nil {
		return err
	}

	return controller.messageInfoDao.RemoveComment(authorUUID, messageId)
}

func (controller MessageControllerImpl) ForwardMessage(conversationId int64, messageId int64, authorUUID string, forwardToConversationID int64) (views.MessageView, error) {
	// Check if the author is in both conversations
	if res, err := controller.conversationDao.IsParticipant(conversationId, authorUUID); err != nil {
		return views.MessageView{}, err
	} else if !res {
		return views.MessageView{}, NotParticipant
	}

	if res, err := controller.conversationDao.IsParticipant(forwardToConversationID, authorUUID); err != nil {
		return views.MessageView{}, err
	} else if !res {
		return views.MessageView{}, NotParticipant
	}

	// Get the message to forward
	message, err := controller.messageDao.GetConversationMessage(messageId, conversationId)
	if err != nil {
		return views.MessageView{}, err
	}

	var attachmentUUID *string = nil
	if attachment := message.GetAttachment(); attachment != nil {
		tmpAttachmentUUID := attachment.GetUUID()
		attachmentUUID = &tmpAttachmentUUID
	}

	// Create the new message copying the previous
	newMessage, err := controller.messageDao.CreateMessage(
		forwardToConversationID, authorUUID, nil, message.GetContent(), attachmentUUID,
	)
	if err != nil {
		return views.MessageView{}, err
	}

	return controller.GetConversationMessage(conversationId, newMessage.GetId(), authorUUID)
}

func (controller MessageControllerImpl) SetAllMessageDelivered(requestIssuerUUID string) ([]database.UserConversation, error) {
	if err := controller.messageDao.SetMessagesAsDelivered(requestIssuerUUID); err != nil {
		return []database.UserConversation{}, err
	}

	return controller.conversationController.GetUserConversations(requestIssuerUUID)
}

func (controller MessageControllerImpl) SetConversationMessagesAsSeen(conversationId int64, requestIssuer string) ([]views.MessageView, error) {
	if err := controller.messageDao.SetConversationMessagesAsSeen(conversationId, requestIssuer); err != nil {
		return nil, err
	}

	return controller.GetConversationMessages(conversationId, 0, 20, requestIssuer)
}
