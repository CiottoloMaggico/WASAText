package models

import (
	"github.com/ciottolomaggico/wasatext/service/database"
)

type MessageInfo struct {
	Message int64   `db:"messageInfo_message"`
	User    string  `db:"messageInfo_user"`
	Status  int     `db:"messageInfo_status"`
	Comment *string `db:"messageInfo_comment"`
}

type MessageInfoWithUser struct {
	Message int64   `db:"messageInfo_message"`
	Status  int     `db:"messageInfo_status"`
	Comment *string `db:"messageInfo_comment"`
	UserWithImage
}

type MessageInfoModel interface {
	GetMessageInfo(user string, message int64) (*MessageInfo, error)
	GetMessageComments(messageId int64) ([]MessageInfoWithUser, error)
	SetComment(user string, message int64, comment string) (*MessageInfo, error)
	RemoveComment(user string, message int64) error
	CountMessageComments(messageId int64) (int, error)
}

type MessageInfoModelImpl struct {
	Db database.AppDatabase
}

func (model MessageInfoModelImpl) CountMessageComments(messageId int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM User_Message WHERE message = ? AND comment IS NOT NULL;`

	if err := model.Db.QueryRow(query, messageId).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (model MessageInfoModelImpl) GetMessageComments(messageId int64) ([]MessageInfoWithUser, error) {
	query := `
		SELECT um.message messageInfo_message, um.status messageInfo_status, um.comment messageInfo_comment, vu.*
		FROM User_Message um, ViewUsers vu
		WHERE message = ? AND vu.user_uuid = um.user AND comment IS NOT NULL;
	`

	messageInfos := make([]MessageInfoWithUser, 0)
	if err := model.Db.QueryStruct(&messageInfos, query, messageId); err != nil {
		return nil, err
	}
	return messageInfos, nil
}

func (model MessageInfoModelImpl) GetMessageInfo(user string, message int64) (*MessageInfo, error) {
	query := `
		SELECT message messageInfo_message, user messageInfo_user, status messageInfo_status, comment messageInfo_comment
		FROM User_Message
		WHERE user = ? AND message = ?
		`

	messageInfo := MessageInfo{}
	if err := model.Db.QueryStructRow(&messageInfo, query, user, message); err != nil {
		return nil, err
	}
	return &messageInfo, nil
}

func (model MessageInfoModelImpl) SetComment(user string, message int64, comment string) (*MessageInfo, error) {
	query := `
		UPDATE User_Message
		SET comment = ?
		WHERE user = ? AND message = ?
		RETURNING message messageInfo_message, user messageInfo_user, status messageInfo_status, comment messageInfo_comment;
	`

	messageInfo := MessageInfo{}
	if err := model.Db.QueryStructRow(&messageInfo, query, comment, user, message); err != nil {
		return nil, err
	}
	return &messageInfo, nil
}

func (model MessageInfoModelImpl) RemoveComment(user string, message int64) error {
	query := `
		UPDATE User_Message
		SET comment = NULL
		WHERE user = ? AND message = ?
		RETURNING message messageInfo_message, user messageInfo_user, status messageInfo_status, comment messageInfo_comment;
`

	_, err := model.Db.Exec(query, user, message)
	return err
}
