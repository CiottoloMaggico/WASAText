package controllers

import (
	"errors"
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
	"io"
)

// TODO: handle pagination

var NotParticipant = errors.New("this user isn't a participant of this group")

type ConversationController interface {
	CreateGroup(groupName string, authorUUID string, photoExtension *string, photoFile *io.Reader) (views.UserConversationView, error)
	LeaveGroup(groupId int64, requestIssuerUUID string) error
	AddToGroup(groupId int64, requestIssuerUUID string, newParticipants []string) (views.UserConversationView, error)
	ChangeGroupName(groupId int64, requestIssuerUUID string, newName string) (views.UserConversationView, error)
	ChangeGroupPhoto(groupId int64, requestIssuerUUID string, photoExtension string, photoFile io.Reader) (views.UserConversationView, error)
	CreateChat(authorUUID string, recipientUUID string) (views.UserConversationView, error)
	GetUserConversations(requestIssuerUUID string, page uint, pageSize uint) ([]views.UserConversationView, error)
	GetUserConversation(requestIssuerUUID string, conversationId int64) (views.UserConversationView, error)
}

type ConversationControllerImpl struct {
	ImageController       ImageController
	ConversationModel     models.ConversationModel
	UserConversationModel models.UserConversationModel
}

func (controller ConversationControllerImpl) CreateGroup(groupName string, authorUUID string, photoExtension *string, photoFile *io.Reader) (views.UserConversationView, error) {
	var photoUUID *string = nil

	if photoExtension != nil && photoFile != nil {
		photo, err := controller.ImageController.CreateImage(*photoExtension, *photoFile)
		if err != nil {
			return views.UserConversationView{}, err
		}
		photoUUID = &photo.Uuid
	}

	commit := false
	defer func(commit bool) {
		if !commit && photoUUID != nil {
			controller.ImageController.DeleteImage(*photoUUID)
		}
	}(commit)

	group, err := controller.ConversationModel.CreateGroup(groupName, authorUUID, photoUUID)
	if err != nil {
		return views.UserConversationView{}, err
	}

	commit = true
	return controller.GetUserConversation(authorUUID, group.Id)
}

func (controller ConversationControllerImpl) LeaveGroup(groupId int64, requestIssuerUUID string) error {
	if res, err := controller.ConversationModel.IsParticipant(groupId, requestIssuerUUID); err != nil {
		return err
	} else if !res {
		return NotParticipant
	}

	return controller.ConversationModel.RemoveGroupParticipant(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) AddToGroup(groupId int64, requestIssuerUUID string, newParticipants []string) (views.UserConversationView, error) {
	if res, err := controller.ConversationModel.IsParticipant(groupId, requestIssuerUUID); err != nil {
		return views.UserConversationView{}, err
	} else if !res {
		return views.UserConversationView{}, NotParticipant
	}

	if err := controller.ConversationModel.AddGroupParticipants(newParticipants, groupId); err != nil {
		return views.UserConversationView{}, err
	}

	return controller.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) ChangeGroupName(groupId int64, requestIssuerUUID string, newName string) (views.UserConversationView, error) {
	if res, err := controller.ConversationModel.IsParticipant(groupId, requestIssuerUUID); err != nil {
		return views.UserConversationView{}, err
	} else if !res {
		return views.UserConversationView{}, NotParticipant
	}

	if _, err := controller.ConversationModel.UpdateGroupName(groupId, newName); err != nil {
		return views.UserConversationView{}, err
	}

	return controller.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) ChangeGroupPhoto(groupId int64, requestIssuerUUID string, photoExtension string, photoFile io.Reader) (views.UserConversationView, error) {
	if res, err := controller.ConversationModel.IsParticipant(groupId, requestIssuerUUID); err != nil {
		return views.UserConversationView{}, err
	} else if !res {
		return views.UserConversationView{}, NotParticipant
	}

	image, err := controller.ImageController.CreateImage(photoExtension, photoFile)
	if err != nil {
		return views.UserConversationView{}, err
	}

	if _, err := controller.ConversationModel.UpdateGroupPic(groupId, image.Uuid); err != nil {
		return views.UserConversationView{}, err
	}

	return controller.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) CreateChat(authorUUID string, recipientUUID string) (views.UserConversationView, error) {
	chat, err := controller.ConversationModel.CreateChat(authorUUID, recipientUUID)
	if err != nil {
		return views.UserConversationView{}, err
	}

	return controller.GetUserConversation(authorUUID, chat.Id)
}

func (controller ConversationControllerImpl) GetUserConversations(requestIssuerUUID string, page uint, pageSize uint) ([]views.UserConversationView, error) {
	conversations, err := controller.UserConversationModel.GetUserConversations(requestIssuerUUID, page, pageSize)
	if errors.Is(err, database.NoResult) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return translators.UserConversationListToView(conversations), nil
}

func (controller ConversationControllerImpl) GetUserConversation(requestIssuerUUID string, conversationId int64) (views.UserConversationView, error) {
	conversation, err := controller.UserConversationModel.GetUserConversation(requestIssuerUUID, conversationId)
	if err != nil {
		return views.UserConversationView{}, err
	}

	return translators.UserConversationToView(*conversation), nil
}
