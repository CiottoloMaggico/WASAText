package controllers

import (
	"errors"
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/ciottolomaggico/wasatext/service/views/pagination"
	"io"
)

type ConversationController interface {
	CreateGroup(groupName string, authorUUID string, photoExtension *string, photoFile io.ReadSeeker) (views.UserConversationView, error)
	LeaveGroup(groupId int64, requestIssuerUUID string) error
	AddToGroup(groupId int64, requestIssuerUUID string, newParticipants []string) (views.UserConversationView, error)
	ChangeGroupName(groupId int64, requestIssuerUUID string, newName string) (views.UserConversationView, error)
	ChangeGroupPhoto(groupId int64, requestIssuerUUID string, photoExtension string, photoFile io.ReadSeeker) (views.UserConversationView, error)
	CreateChat(authorUUID string, recipientUUID string) (views.UserConversationView, error)
	GetUserConversations(requestIssuerUUID string, paginationPs pagination.PaginationParams) (pagination.PaginatedView, error)
	GetUserConversation(requestIssuerUUID string, conversationId int64) (views.UserConversationView, error)
	IsParticipant(conversationId int64, userUUID string) (bool, error)
}

type ConversationControllerImpl struct {
	ImageController       ImageController
	ConversationModel     models.ConversationModel
	UserConversationModel models.UserConversationModel
}

func (controller ConversationControllerImpl) IsParticipant(conversationId int64, userUUID string) (bool, error) {
	if ok, err := controller.ConversationModel.IsParticipant(conversationId, userUUID); !ok {
		return false, api_errors.Forbidden()
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (controller ConversationControllerImpl) CreateGroup(groupName string, authorUUID string, photoExtension *string, photoFile io.ReadSeeker) (views.UserConversationView, error) {
	var photoUUID *string = nil

	if photoExtension != nil && photoFile != nil {
		photo, err := controller.ImageController.CreateImage(*photoExtension, photoFile)
		if err != nil {
			return views.UserConversationView{}, err
		}
		photoUUID = &photo.Uuid
	}

	group, err := controller.ConversationModel.CreateGroup(groupName, authorUUID, photoUUID)
	if err != nil {
		return views.UserConversationView{}, err
	}

	return controller.GetUserConversation(authorUUID, group.Id)
}

func (controller ConversationControllerImpl) LeaveGroup(groupId int64, requestIssuerUUID string) error {
	if conv, err := controller.GetUserConversation(requestIssuerUUID, groupId); conv.Type != "group" {
		return api_errors.ResourceNotFound()
	} else if err != nil {
		return err
	}

	return controller.ConversationModel.RemoveGroupParticipant(requestIssuerUUID, groupId)

}

func (controller ConversationControllerImpl) AddToGroup(groupId int64, requestIssuerUUID string, newParticipants []string) (views.UserConversationView, error) {
	if conv, err := controller.GetUserConversation(requestIssuerUUID, groupId); conv.Type != "group" {
		return views.UserConversationView{}, api_errors.ResourceNotFound()
	} else if err != nil {
		return views.UserConversationView{}, err
	}

	if err := controller.ConversationModel.AddGroupParticipants(newParticipants, groupId); errors.Is(err, database.TriggerConstraint) {
		return views.UserConversationView{}, api_errors.Conflict(map[string]string{"group capacity": "Groups can handle al most 200 participants"})
	} else if err != nil {
		return views.UserConversationView{}, err
	}

	return controller.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) ChangeGroupName(groupId int64, requestIssuerUUID string, newName string) (views.UserConversationView, error) {
	if conv, err := controller.GetUserConversation(requestIssuerUUID, groupId); conv.Type != "group" {
		return views.UserConversationView{}, api_errors.ResourceNotFound()
	} else if err != nil {
		return views.UserConversationView{}, err
	}

	if _, err := controller.ConversationModel.UpdateGroupName(groupId, newName); err != nil {
		return views.UserConversationView{}, err
	}

	return controller.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) ChangeGroupPhoto(groupId int64, requestIssuerUUID string, photoExtension string, photoFile io.ReadSeeker) (views.UserConversationView, error) {
	if conv, err := controller.GetUserConversation(requestIssuerUUID, groupId); conv.Type != "group" {
		return views.UserConversationView{}, api_errors.ResourceNotFound()
	} else if err != nil {
		return views.UserConversationView{}, err
	}

	image, err := controller.ImageController.CreateImage(photoExtension, photoFile)
	if err != nil {
		return views.UserConversationView{}, err
	}

	if _, err = controller.ConversationModel.UpdateGroupPic(groupId, image.Uuid); err != nil {
		if err = controller.ImageController.DeleteImage(image.Uuid); err != nil {
			return views.UserConversationView{}, err
		}
		return views.UserConversationView{}, err
	}

	return controller.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) CreateChat(authorUUID string, recipientUUID string) (views.UserConversationView, error) {
	chat, err := controller.ConversationModel.CreateChat(authorUUID, recipientUUID)
	if errors.Is(err, database.UniqueConstraint) {
		return views.UserConversationView{}, api_errors.NewApiError(409, "The chat already exists")
	} else if err != nil {
		return views.UserConversationView{}, err
	}

	return controller.GetUserConversation(authorUUID, chat.Id)
}

func (controller ConversationControllerImpl) GetUserConversations(requestIssuerUUID string, paginationPs pagination.PaginationParams) (pagination.PaginatedView, error) {
	conversationsCount, err := controller.UserConversationModel.Count(requestIssuerUUID)
	if err != nil {
		return pagination.PaginatedView{}, err
	}

	if conversationsCount == 0 {
		return pagination.ToPaginatedView(paginationPs, conversationsCount, translators.UserConversationListToSummaryView(make([]models.UserConversation, 0)))
	}

	conversations, err := controller.UserConversationModel.GetUserConversations(requestIssuerUUID, paginationPs.Page, paginationPs.Size)
	if err != nil {
		return pagination.PaginatedView{}, nil
	}

	return pagination.ToPaginatedView(paginationPs, conversationsCount, translators.UserConversationListToSummaryView(conversations))
}

func (controller ConversationControllerImpl) GetUserConversation(requestIssuerUUID string, conversationId int64) (views.UserConversationView, error) {
	conversation, err := controller.UserConversationModel.GetUserConversation(requestIssuerUUID, conversationId)
	if errors.Is(err, database.NoResult) {
		return views.UserConversationView{}, api_errors.ResourceNotFound()
	} else if err != nil {
		return views.UserConversationView{}, err
	}

	participants, err := controller.ConversationModel.GetConversationParticipants(conversationId)
	if err != nil {
		return views.UserConversationView{}, err
	}

	return translators.UserConversationToView(*conversation, participants), nil
}
