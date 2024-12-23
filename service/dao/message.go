package dao

import (
	"github.com/ciottolomaggico/wasatext/service/database"
)

type MessageDao interface {
	CreateMessage(
		conversation int64,
		author string,
		replyTo *int64,
		content *string,
		attachment *string,
	) (*database.Message, error)
	DeleteMessage(id int64) error
	GetConversationMessage(messageId int64, conversationId int64) (*database.MessageWithAuthorAndAttachment, error)
	GetConversationMessages(conversationId int64, page uint, size uint) (database.MessageWithAuthorAndAttachmentList, error)
	SetMessagesAsDelivered(user string) error
	SetConversationMessagesAsSeen(conversationId int64, user string) error
}

type MessageDaoImpl struct {
	db database.BaseDatabase
}

func (dao MessageDaoImpl) CreateMessage(
	conversation int64,
	author string,
	replyTo *int64,
	content *string,
	attachment *string,
) (*database.Message, error) {
	query := `INSERT INTO Message (conversation, author, replyTo, content, attachment) VALUES (?, ?, ?, ?, ?) RETURNING *`

	return dao.db.QueryRowMessage(query, conversation, author, replyTo, content, attachment)
}

func (dao MessageDaoImpl) DeleteMessage(id int64) error {
	if _, err := dao.db.Exec(
		`DELETE FROM Message WHERE id = ?;`,
		id,
	); err != nil {
		return err
	}
	return nil
}

func (dao MessageDaoImpl) GetConversationMessage(messageId int64, conversationId int64) (*database.MessageWithAuthorAndAttachment, error) {
	query := `SELECT * FROM ViewMessages WHERE id = ? AND conversation = ?;`

	return dao.db.QueryRowMessageWithAuthorAndAttachment(query, messageId, conversationId)
}

func (dao MessageDaoImpl) GetConversationMessages(conversationId int64, page uint, size uint) (database.MessageWithAuthorAndAttachmentList, error) {
	query := `SELECT * FROM ViewMessages WHERE conversation = ? ORDER BY sendAt DESC LIMIT ? OFFSET ?;`

	return dao.db.QueryMessageWithAuthorAndAttachment(query, conversationId, page, page*size)
}

func (dao MessageDaoImpl) SetMessagesAsDelivered(user string) error {
	query := `UPDATE User_Message SET status = 2 WHERE user = ? AND status = 1 RETURNING *;`
	_, err := dao.db.Exec(query, user)
	return err
}

func (dao MessageDaoImpl) SetConversationMessagesAsSeen(conversationId int64, user string) error {
	query := `
				UPDATE um
				SET status = 3
				FROM Message m, User_Message um
				WHERE um.user = ? AND m.conversation = ?
				  	AND um.message = m.id AND um.status < 3
				RETURNING *;
	`
	_, err := dao.db.Exec(query, user)
	return err
}
