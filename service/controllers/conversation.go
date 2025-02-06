package controllers

import (
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
	"io"
)

type ConversationController interface {
	CreateGroup(groupName string, authorUUID string, photoExtension *string, photoFile io.ReadSeeker) (views.UserConversationView, error)
	LeaveGroup(groupId int64, requestIssuerUUID string) error
	AddToGroup(groupId int64, requestIssuerUUID string, newParticipants []string) (views.UserConversationView, error)
	ChangeGroupName(groupId int64, requestIssuerUUID string, newName string) (views.UserConversationView, error)
	ChangeGroupPhoto(groupId int64, requestIssuerUUID string, photoExtension string, photoFile io.ReadSeeker) (views.UserConversationView, error)
	CreateChat(authorUUID string, recipientUUID string) (views.UserConversationView, error)
	IsParticipant(conversationId int64, userUUID string) (bool, error)
}

type ConversationControllerImpl struct {
	Model                      models.ConversationModel
	ImageController            ImageController
	UserConversationController UserConversationController
}

func (controller ConversationControllerImpl) IsParticipant(conversationId int64, userUUID string) (bool, error) {
	if ok, err := controller.Model.IsParticipant(conversationId, userUUID); !ok {
		return false, api_errors.Forbidden()
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (controller ConversationControllerImpl) CreateGroup(groupName string, authorUUID string, photoExtension *string, photoFile io.ReadSeeker) (views.UserConversationView, error) {
	var photoUUID *string = nil

	commit := true
	defer func() {
		if !commit {
			_ = controller.ImageController.DeleteImage(*photoUUID)
		}
	}()

	if photoExtension != nil && photoFile != nil {
		photo, err := controller.ImageController.CreateImage(*photoExtension, photoFile)
		if err != nil {
			return views.UserConversationView{}, err
		}
		photoUUID = &photo.Uuid
		commit = false
	}

	group, err := controller.Model.CreateGroup(groupName, authorUUID, photoUUID)
	if err != nil {
		return views.UserConversationView{}, err
	}

	commit = true

	return controller.UserConversationController.GetUserConversation(authorUUID, group.Id)
}

func (controller ConversationControllerImpl) LeaveGroup(groupId int64, requestIssuerUUID string) error {
	if conv, err := controller.UserConversationController.GetUserConversation(requestIssuerUUID, groupId); conv.Type != "group" {
		return api_errors.ResourceNotFound()
	} else if err != nil {
		return err
	}

	return controller.Model.RemoveGroupParticipant(requestIssuerUUID, groupId)

}

func (controller ConversationControllerImpl) AddToGroup(groupId int64, requestIssuerUUID string, newParticipants []string) (views.UserConversationView, error) {
	if conv, err := controller.UserConversationController.GetUserConversation(requestIssuerUUID, groupId); conv.Type != "group" {
		return views.UserConversationView{}, api_errors.ResourceNotFound()
	} else if err != nil {
		return views.UserConversationView{}, err
	}

	if err := controller.Model.AddGroupParticipants(newParticipants, groupId); err != nil {
		return views.UserConversationView{}, translators.DBErrorToApiError(err)
	}

	return controller.UserConversationController.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) ChangeGroupName(groupId int64, requestIssuerUUID string, newName string) (views.UserConversationView, error) {
	if conv, err := controller.UserConversationController.GetUserConversation(requestIssuerUUID, groupId); conv.Type != "group" {
		return views.UserConversationView{}, api_errors.ResourceNotFound()
	} else if err != nil {
		return views.UserConversationView{}, err
	}

	if _, err := controller.Model.UpdateGroupName(groupId, newName); err != nil {
		return views.UserConversationView{}, err
	}

	return controller.UserConversationController.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) ChangeGroupPhoto(groupId int64, requestIssuerUUID string, photoExtension string, photoFile io.ReadSeeker) (views.UserConversationView, error) {
	if conv, err := controller.UserConversationController.GetUserConversation(requestIssuerUUID, groupId); conv.Type != "group" {
		return views.UserConversationView{}, api_errors.ResourceNotFound()
	} else if err != nil {
		return views.UserConversationView{}, err
	}

	image, err := controller.ImageController.CreateImage(photoExtension, photoFile)
	if err != nil {
		return views.UserConversationView{}, err
	}

	commit := false
	defer func() {
		if !commit {
			_ = controller.ImageController.DeleteImage(image.Uuid)
		}
	}()

	if _, err = controller.Model.UpdateGroupPic(groupId, image.Uuid); err != nil {
		return views.UserConversationView{}, err
	}

	commit = true

	return controller.UserConversationController.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) CreateChat(authorUUID string, recipientUUID string) (views.UserConversationView, error) {
	chat, err := controller.Model.CreateChat(authorUUID, recipientUUID)
	if err != nil {
		return views.UserConversationView{}, translators.DBErrorToApiError(err)
	}

	return controller.UserConversationController.GetUserConversation(authorUUID, chat.Id)
}
