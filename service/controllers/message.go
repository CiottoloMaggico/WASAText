package controllers

import (
	apierrors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/filter"
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/ciottolomaggico/wasatext/service/views/pagination"
	"io"
)

type MessageController interface {
	SendMessage(conversationID int64, authorUUID string, replyToId *int64, content *string, attachmentExt *string, attachmentFile io.ReadSeeker) (views.MessageView, error)
	DeleteMessage(conversationID int64, messageId int64, requestIssuerUUID string) error
	GetConversationMessages(conversationID int64, requestIssuerUUID string, paginationPs pagination.PaginationParams) (pagination.PaginatedView, error)
	GetConversationMessage(conversationID int64, messageId int64, requestIssuerUUID string) (views.MessageView, error)
	ForwardMessage(conversationId int64, messageId int64, authorUUID string, forwardToConversationID int64) (views.MessageView, error)
	SetAllMessageDelivered(requestIssuerUUID string, paginationPs pagination.PaginationParams) (pagination.PaginatedView, error)
	SetConversationMessagesAsSeen(conversationId int64, requestIssuer string, paginationPs pagination.PaginationParams) (pagination.PaginatedView, error)
}

type MessageControllerImpl struct {
	Model                      models.MessageModel
	ConversationController     ConversationController
	UserConversationController UserConversationController
	ImageController            ImageController
	Filter                     filter.Filter
}

func (controller MessageControllerImpl) SendMessage(conversationID int64, authorUUID string, replyToId *int64, content *string, attachmentExt *string, attachmentFile io.ReadSeeker) (views.MessageView, error) {
	if _, err := controller.ConversationController.IsParticipant(conversationID, authorUUID); err != nil {
		return views.MessageView{}, err
	}

	var attachmentUUID *string = nil

	commit := true
	defer func() {
		if !commit {
			_ = controller.ImageController.DeleteImage(*attachmentUUID)
		}
	}()

	if attachmentExt != nil && attachmentFile != nil {
		image, err := controller.ImageController.CreateImage(*attachmentExt, attachmentFile)
		if err != nil {
			return views.MessageView{}, err
		}
		attachmentUUID = &image.Uuid
		commit = false
	}

	message, err := controller.Model.CreateMessage(
		conversationID, authorUUID, replyToId, content, attachmentUUID,
	)
	if err != nil {
		return views.MessageView{}, translators.DBErrorToApiError(err)
	}

	commit = true
	return controller.GetConversationMessage(conversationID, message.Id, authorUUID)
}

func (controller MessageControllerImpl) GetConversationMessage(conversationID int64, messageId int64, requestIssuerUUID string) (views.MessageView, error) {
	if _, err := controller.ConversationController.IsParticipant(conversationID, requestIssuerUUID); err != nil {
		return views.MessageView{}, err
	}

	message, err := controller.Model.GetConversationMessage(messageId, conversationID)
	if err != nil {
		return views.MessageView{}, translators.DBErrorToApiError(err)
	}

	return translators.MessageWithAuthorAndAttachmentToView(*message), err
}

func (controller MessageControllerImpl) GetConversationMessages(conversationID int64, requestIssuerUUID string, paginationPs pagination.PaginationParams) (pagination.PaginatedView, error) {
	if _, err := controller.ConversationController.IsParticipant(conversationID, requestIssuerUUID); err != nil {
		return pagination.PaginatedView{}, err
	}

	filterQuery, err := controller.Filter.Evaluate(paginationPs.Filter)
	if err != nil {
		return pagination.PaginatedView{}, apierrors.InvalidUrlParameters()
	}
	queryParameters := database.NewQueryParameters(paginationPs.Page, paginationPs.Size, filterQuery)

	messagesCount, err := controller.Model.Count(conversationID, queryParameters)
	if err != nil {
		return pagination.PaginatedView{}, err
	}

	messages := make([]models.MessageWithAuthorAndAttachment, 0)
	if messagesCount > 0 {
		messages, err = controller.Model.GetConversationMessages(conversationID, queryParameters)
		if err != nil {
			return pagination.PaginatedView{}, err
		}
	}

	return pagination.ToPaginatedView(paginationPs, messagesCount, translators.MessageWithAuthorAndAttachmentListToView(messages))
}

func (controller MessageControllerImpl) DeleteMessage(conversationID int64, messageId int64, requestIssuerUUID string) error {
	message, err := controller.GetConversationMessage(conversationID, messageId, requestIssuerUUID)
	if err != nil {
		return err
	}

	if message.Author.Uuid != requestIssuerUUID {
		return apierrors.Forbidden()
	}

	if err := controller.Model.DeleteMessage(message.Id); err != nil {
		return err
	}

	return nil
}

func (controller MessageControllerImpl) ForwardMessage(conversationId int64, messageId int64, authorUUID string, forwardToConversationID int64) (views.MessageView, error) {
	// check if the requestIssuer ("authorUUID") is a participant of the dest conversation
	if _, err := controller.ConversationController.IsParticipant(forwardToConversationID, authorUUID); err != nil {
		return views.MessageView{}, err
	}

	// Get the message to forward
	message, err := controller.GetConversationMessage(conversationId, messageId, authorUUID)
	if err != nil {
		return views.MessageView{}, err
	}

	var attachmentUUID *string = nil
	if message.Attachment != nil {
		attachmentUUID = &message.Attachment.Uuid
	}

	// Create the new message copying the previous
	newMessage, err := controller.Model.CreateMessage(
		forwardToConversationID, authorUUID, nil, message.Content, attachmentUUID,
	)
	if err != nil {
		return views.MessageView{}, err
	}

	return controller.GetConversationMessage(forwardToConversationID, newMessage.Id, authorUUID)
}

func (controller MessageControllerImpl) SetAllMessageDelivered(requestIssuerUUID string, paginationPs pagination.PaginationParams) (pagination.PaginatedView, error) {
	if err := controller.Model.SetMessagesAsDelivered(requestIssuerUUID); err != nil {
		return pagination.PaginatedView{}, err
	}

	return controller.UserConversationController.GetUserConversations(requestIssuerUUID, paginationPs)
}

func (controller MessageControllerImpl) SetConversationMessagesAsSeen(conversationId int64, requestIssuer string, paginationPs pagination.PaginationParams) (pagination.PaginatedView, error) {
	if err := controller.Model.SetConversationMessagesAsSeen(conversationId, requestIssuer); err != nil {
		return pagination.PaginatedView{}, err
	}

	return controller.GetConversationMessages(conversationId, requestIssuer, paginationPs)
}
