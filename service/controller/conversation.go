package controllers

import (
	"database/sql"
	"errors"
	"github.com/ciottolomaggico/wasatext/service/dao"
	"github.com/ciottolomaggico/wasatext/service/database"
	"io"
)

// TODO: handle pagination

var NotParticipant = errors.New("this user isn't a participant of this group")

type ConversationController interface {
	CreateGroup(groupName string, authorUUID string, photoUUID *string) (database.UserConversation, error)
	LeaveGroup(groupId int64, requestIssuerUUID string) error
	AddToGroup(groupId int64, requestIssuerUUID string, newParticipantUUID string) (database.UserConversation, error)
	ChangeGroupName(groupId int64, requestIssuerUUID string, newName string) (database.UserConversation, error)
	ChangeGroupPhoto(groupId int64, requestIssuerUUID string, newPhotoUUID string) (database.UserConversation, error)
	CreateChat(authorUUID string, recipientUUID string) (database.UserConversation, error)
	GetUserConversations(requestIssuerUUID string) ([]database.UserConversation, error)
	GetUserConversation(requestIssuerUUID string, conversationId int64) (database.UserConversation, error)
}

type ConversationControllerImpl struct {
	imageController     ImageController
	conversationDao     dao.ConversationDao
	userConversationDao dao.UserConversationDao
}

func (controller ConversationControllerImpl) CreateGroup(groupName string, authorUUID string, photoUUID *string) (database.UserConversation, error) {
	group, err := controller.conversationDao.CreateGroup(groupName, authorUUID, photoUUID)
	if err != nil {
		return database.UserConversation{}, err
	}

	return controller.GetUserConversation(authorUUID, group.GetId())
}

func (controller ConversationControllerImpl) LeaveGroup(groupId int64, requestIssuerUUID string) error {
	if res, err := controller.conversationDao.IsParticipant(groupId, requestIssuerUUID); err != nil {
		return err
	} else if !res {
		return NotParticipant
	}

	return controller.conversationDao.RemoveGroupParticipant(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) AddToGroup(groupId int64, requestIssuerUUID string, newParticipantUUID string) (database.UserConversation, error) {
	if res, err := controller.conversationDao.IsParticipant(groupId, requestIssuerUUID); err != nil {
		return database.UserConversation{}, err
	} else if !res {
		return database.UserConversation{}, NotParticipant
	}

	if err := controller.conversationDao.AddGroupParticipant(newParticipantUUID, groupId); err != nil {
		return database.UserConversation{}, err
	}

	return controller.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) ChangeGroupName(groupId int64, requestIssuerUUID string, newName string) (database.UserConversation, error) {
	if res, err := controller.conversationDao.IsParticipant(groupId, requestIssuerUUID); err != nil {
		return database.UserConversation{}, err
	} else if !res {
		return database.UserConversation{}, NotParticipant
	}

	if _, err := controller.conversationDao.UpdateGroupName(groupId, newName); err != nil {
		return database.UserConversation{}, err
	}

	return controller.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) ChangeGroupPhoto(groupId int64, requestIssuerUUID string, photoExtension string, photoFile io.Reader) (database.UserConversation, error) {
	if res, err := controller.conversationDao.IsParticipant(groupId, requestIssuerUUID); err != nil {
		return database.UserConversation{}, err
	} else if !res {
		return database.UserConversation{}, NotParticipant
	}

	image, err := controller.imageController.CreateImage(photoExtension, photoFile)
	if err != nil {
		return database.UserConversation{}, err
	}

	if _, err := controller.conversationDao.UpdateGroupPic(groupId, image.Uuid); err != nil {
		return database.UserConversation{}, err
	}

	return controller.GetUserConversation(requestIssuerUUID, groupId)
}

func (controller ConversationControllerImpl) CreateChat(authorUUID string, recipientUUID string) (database.UserConversation, error) {
	chat, err := controller.conversationDao.CreateChat(authorUUID, recipientUUID)
	if err != nil {
		return database.UserConversation{}, err
	}

	return controller.GetUserConversation(authorUUID, chat.GetId())
}

func (controller ConversationControllerImpl) GetUserConversations(requestIssuerUUID string, page uint, pageSize uint) ([]database.UserConversation, error) {
	conversations, err := controller.userConversationDao.GetUserConversations(requestIssuerUUID, page, pageSize)
	if errors.Is(err, sql.ErrNoRows) {
		return conversations, nil
	} else if err != nil {
		return nil, err
	}

	return conversations, nil
}

func (controller ConversationControllerImpl) GetUserConversation(requestIssuerUUID string, conversationId int64) (database.UserConversation, error) {
	conversation, err := controller.userConversationDao.GetUserConversation(requestIssuerUUID, conversationId)
	if err != nil {
		return database.UserConversation{}, err
	}

	return *conversation, nil
}
