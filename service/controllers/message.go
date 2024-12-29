package controllers

import (
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
	"io"
)

type MessageController interface {
	SendMessage(conversationID int64, authorUUID string, replyToId *int64, content *string, attachmentExt *string, attachmentFile *io.Reader) (views.MessageView, error)
	DeleteMessage(conversationID int64, messageId int64, requestIssuerUUID string) error
	GetConversationMessages(conversationID int64, page uint, pageSize uint, requestIssuerUUID string) ([]views.MessageView, error)
	GetConversationMessage(conversationID int64, messageId int64, requestIssuerUUID string) (views.MessageView, error)
	GetComments(conversationId int64, messageId int64, requestIssuerUUID string) ([]views.CommentView, error)
	CommentMessage(conversationId int64, messageId int64, authorUUID string, comment string) (views.CommentView, error)
	UncommentMessage(conversationId int64, messageId int64, authorUUID string) error
	ForwardMessage(conversationId int64, messageId int64, authorUUID string, forwardToConversationID int64) (views.MessageView, error)
	SetAllMessageDelivered(requestIssuerUUID string) ([]views.UserConversationView, error)
	SetConversationMessagesAsSeen(conversationId int64, requestIssuer string) ([]views.MessageView, error)
}

type MessageControllerImpl struct {
	ConversationController ConversationController
	ImageController        ImageController
	MessageInfoModel       models.MessageInfoModel
	MessageModel           models.MessageModel
	ConversationModel      models.ConversationModel
}

func (controller MessageControllerImpl) SendMessage(conversationID int64, authorUUID string, replyToId *int64, content *string, attachmentExt *string, attachmentFile *io.Reader) (views.MessageView, error) {
	if res, err := controller.ConversationModel.IsParticipant(conversationID, authorUUID); err != nil {
		return views.MessageView{}, err
	} else if !res {
		return views.MessageView{}, NotParticipant
	}

	var attachmentUUID *string = nil
	if attachmentExt != nil && attachmentFile != nil {
		image, err := controller.ImageController.CreateImage(*attachmentExt, *attachmentFile)
		if err != nil {
			return views.MessageView{}, err
		}
		attachmentUUID = &image.Uuid
	}

	message, err := controller.MessageModel.CreateMessage(
		conversationID, authorUUID, replyToId, content, attachmentUUID,
	)
	if err != nil {
		return views.MessageView{}, err
	}
	return controller.GetConversationMessage(conversationID, message.Id, authorUUID)
}

func (controller MessageControllerImpl) GetConversationMessage(conversationID int64, messageId int64, requestIssuerUUID string) (views.MessageView, error) {
	if res, err := controller.ConversationModel.IsParticipant(conversationID, requestIssuerUUID); err != nil {
		return views.MessageView{}, err
	} else if !res {
		return views.MessageView{}, NotParticipant
	}

	message, err := controller.MessageModel.GetConversationMessage(messageId, conversationID)
	if err != nil {
		return views.MessageView{}, err
	}
	return translators.MessageWithAuthorAndAttachmentToView(*message), err
}

func (controller MessageControllerImpl) GetConversationMessages(conversationID int64, page uint, pageSize uint, requestIssuerUUID string) ([]views.MessageView, error) {
	if res, err := controller.ConversationModel.IsParticipant(conversationID, requestIssuerUUID); err != nil {
		return []views.MessageView{}, err
	} else if !res {
		return []views.MessageView{}, NotParticipant
	}

	messages, err := controller.MessageModel.GetConversationMessages(conversationID, page, pageSize)
	if err != nil {
		return []views.MessageView{}, err
	}

	return translators.MessageWithAuthorAndAttachmentListToView(messages), nil
}

func (controller MessageControllerImpl) DeleteMessage(conversationID int64, messageId int64, requestIssuerUUID string) error {
	if res, err := controller.ConversationModel.IsParticipant(conversationID, requestIssuerUUID); err != nil {
		return err
	} else if !res {
		return NotParticipant
	}

	_, err := controller.MessageModel.GetConversationMessage(conversationID, messageId)
	if err != nil {
		return err
	}

	if err := controller.MessageModel.DeleteMessage(messageId); err != nil {
		return err
	}

	return nil
}

func (controller MessageControllerImpl) GetComments(conversationId int64, messageId int64, requestIssuerUUID string) ([]views.CommentView, error) {
	if _, err := controller.GetConversationMessage(conversationId, messageId, requestIssuerUUID); err != nil {
		return nil, err
	}

	comments, err := controller.MessageInfoModel.GetMessageComments(messageId)

	if err != nil {
		return nil, err
	}

	return translators.MessageInfoListToCommentView(comments), nil
}

func (controller MessageControllerImpl) CommentMessage(conversationId int64, messageId int64, authorUUID string, comment string) (views.CommentView, error) {
	if res, err := controller.ConversationModel.IsParticipant(conversationId, authorUUID); err != nil {
		return views.CommentView{}, err
	} else if !res {
		return views.CommentView{}, NotParticipant
	}

	if _, err := controller.GetConversationMessage(conversationId, messageId, authorUUID); err != nil {
		return views.CommentView{}, err
	}

	res, err := controller.MessageInfoModel.SetComment(authorUUID, messageId, comment)
	if err != nil {
		return views.CommentView{}, err
	}
	return translators.MessageInfoToCommentView(*res), nil
}

func (controller MessageControllerImpl) UncommentMessage(conversationId int64, messageId int64, authorUUID string) error {
	if res, err := controller.ConversationModel.IsParticipant(conversationId, authorUUID); err != nil {
		return err
	} else if !res {
		return NotParticipant
	}

	if _, err := controller.GetConversationMessage(conversationId, messageId, authorUUID); err != nil {
		return err
	}

	return controller.MessageInfoModel.RemoveComment(authorUUID, messageId)
}

func (controller MessageControllerImpl) ForwardMessage(conversationId int64, messageId int64, authorUUID string, forwardToConversationID int64) (views.MessageView, error) {
	// Check if the author is in both conversations
	if res, err := controller.ConversationModel.IsParticipant(conversationId, authorUUID); err != nil {
		return views.MessageView{}, err
	} else if !res {
		return views.MessageView{}, NotParticipant
	}

	if res, err := controller.ConversationModel.IsParticipant(forwardToConversationID, authorUUID); err != nil {
		return views.MessageView{}, err
	} else if !res {
		return views.MessageView{}, NotParticipant
	}

	// Get the message to forward
	message, err := controller.MessageModel.GetConversationMessage(messageId, conversationId)
	if err != nil {
		return views.MessageView{}, err
	}

	// Create the new message copying the previous
	newMessage, err := controller.MessageModel.CreateMessage(
		forwardToConversationID, authorUUID, nil, message.Content, message.AttachmentUuid,
	)
	if err != nil {
		return views.MessageView{}, err
	}

	return controller.GetConversationMessage(conversationId, newMessage.Id, authorUUID)
}

func (controller MessageControllerImpl) SetAllMessageDelivered(requestIssuerUUID string) ([]views.UserConversationView, error) {
	if err := controller.MessageModel.SetMessagesAsDelivered(requestIssuerUUID); err != nil {
		return []views.UserConversationView{}, err
	}

	return controller.ConversationController.GetUserConversations(requestIssuerUUID, 0, 20)
}

func (controller MessageControllerImpl) SetConversationMessagesAsSeen(conversationId int64, requestIssuer string) ([]views.MessageView, error) {
	if err := controller.MessageModel.SetConversationMessagesAsSeen(conversationId, requestIssuer); err != nil {
		return nil, err
	}

	return controller.GetConversationMessages(conversationId, 0, 20, requestIssuer)
}
