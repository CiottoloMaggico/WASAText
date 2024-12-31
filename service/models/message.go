package models

import (
	"github.com/ciottolomaggico/wasatext/service/database"
)

type UserConversationMessagePreview struct {
	Id          *int64  `db:"message_id"`
	SendAt      *string `db:"message_sendAt"`
	DeliveredAt *string `db:"message_deliveredAt"`
	SeenAt      *string `db:"message_seenAt"`
	ReplyTo     *int64  `db:"message_replyTo"`
	Content     *string `db:"message_content"`
	Attachment  *string `db:"message_attachment"`
	Uuid        *string `db:"user_uuid"`
	Username    *string `db:"user_username"`
	Photo       *string `db:"user_image"`
}

type Message struct {
	Id           int64   `db:"message_id"`
	Conversation int64   `db:"message_conversation"`
	Author       string  `db:"message_author"`
	SendAt       string  `db:"message_sendAt"`
	DeliveredAt  *string `db:"message_deliveredAt"`
	SeenAt       *string `db:"message_seenAt"`
	ReplyTo      *int64  `db:"message_replyTo"`
	Content      *string `db:"message_content"`
	Attachment   *string `db:"message_attachment"`
}

type MessageWithAuthor struct {
	Id           int64   `db:"message_id"`
	Conversation int64   `db:"message_conversation"`
	SendAt       string  `db:"message_sendAt"`
	DeliveredAt  *string `db:"message_deliveredAt"`
	SeenAt       *string `db:"message_seenAt"`
	ReplyTo      *int64  `db:"message_replyTo"`
	Content      *string `db:"message_content"`
	Attachment   *string `db:"message_attachment"`
	User
}

type MessageWithAuthorAndAttachment struct {
	Id                   int64   `db:"message_id"`
	Conversation         int64   `db:"message_conversation"`
	SendAt               string  `db:"message_sendAt"`
	DeliveredAt          *string `db:"message_deliveredAt"`
	SeenAt               *string `db:"message_seenAt"`
	ReplyTo              *int64  `db:"message_replyTo"`
	Content              *string `db:"message_content"`
	AttachmentUuid       *string `db:"attachment_uuid"`
	AttachmentExtension  *string `db:"attachment_extension"`
	AttachmentWidth      *int    `db:"attachment_width"`
	AttachmentHeight     *int    `db:"attachment_height"`
	AttachmentFullUrl    *string `db:"attachment_fullUrl"`
	AttachmentUploadedAt *string `db:"attachment_uploadedAt"`
	UserWithImage
}

type MessageModel interface {
	CreateMessage(
		conversation int64,
		author string,
		replyTo *int64,
		content *string,
		attachment *string,
	) (*Message, error)
	DeleteMessage(id int64) error
	GetConversationMessage(messageId int64, conversationId int64) (*MessageWithAuthorAndAttachment, error)
	GetConversationMessages(conversationId int64, page int, size int) ([]MessageWithAuthorAndAttachment, error)
	SetMessagesAsDelivered(user string) error
	SetConversationMessagesAsSeen(conversationId int64, user string) error
	Count(conversationId int64) (int, error)
}

type MessageModelImpl struct {
	Db database.AppDatabase
}

func (m MessageWithAuthorAndAttachment) GetStatus() string {
	res := "sent"

	if m.DeliveredAt != nil {
		res = "delivered"
	}

	if m.SeenAt != nil {
		res = "seen"
	}

	return res
}

func (m MessageWithAuthor) GetStatus() string {
	res := "sent"

	if m.DeliveredAt != nil {
		res = "delivered"
	}

	if m.SeenAt != nil {
		res = "seen"
	}

	return res
}

func (model MessageModelImpl) CreateMessage(
	conversation int64,
	author string,
	replyTo *int64,
	content *string,
	attachment *string,
) (*Message, error) {
	query := `
		INSERT INTO Message (conversation, author, replyTo, content, attachment)
		VALUES (?, ?, ?, ?, ?)
		RETURNING
			id message_id, conversation message_conversation,
			author message_author, sendAt message_sendAt, deliveredAt message_deliveredAt,
			seenAt message_seenAt, replyTo message_replyTo, content message_content, attachment message_attachment;
	`

	message := Message{}
	if err := model.Db.QueryStructRow(&message, query, conversation, author, replyTo, content, attachment); err != nil {
		return nil, err
	}

	return &message, nil
}

func (model MessageModelImpl) DeleteMessage(id int64) error {
	if _, err := model.Db.Exec(
		`DELETE FROM Message WHERE id = ?;`,
		id,
	); err != nil {
		return database.HandleSqlError(err)
	}
	return nil
}

func (model MessageModelImpl) GetConversationMessage(messageId int64, conversationId int64) (*MessageWithAuthorAndAttachment, error) {
	query := `SELECT * FROM ViewMessages WHERE message_id = ? AND message_conversation = ?;`

	message := MessageWithAuthorAndAttachment{}
	if err := model.Db.QueryStructRow(&message, query, messageId, conversationId); err != nil {
		return nil, err
	}

	return &message, nil
}

func (model MessageModelImpl) GetConversationMessages(conversationId int64, page int, size int) ([]MessageWithAuthorAndAttachment, error) {
	query := `SELECT * FROM ViewMessages WHERE message_conversation = ? ORDER BY message_sendAt DESC LIMIT ? OFFSET ?;`

	messages := make([]MessageWithAuthorAndAttachment, 0, size)
	if err := model.Db.QueryStruct(&messages, query, conversationId, page, (page-1)*size); err != nil {
		return nil, err
	}

	return messages, nil
}

func (model MessageModelImpl) SetMessagesAsDelivered(user string) error {
	query := `UPDATE User_Message SET status = 2 WHERE user = ? AND status = 1;`
	if _, err := model.Db.Exec(query, user); err != nil {
		return database.HandleSqlError(err)
	}
	return nil
}

func (model MessageModelImpl) SetConversationMessagesAsSeen(conversationId int64, user string) error {
	query := `
				UPDATE User_Message
				SET status = 3
				FROM Message m, User_Message um
				WHERE um.user = ? AND m.conversation = ?
				  	AND um.message = m.id AND um.status < 3;
	`
	if _, err := model.Db.Exec(query, user, conversationId); err != nil {
		return database.HandleSqlError(err)
	}
	return nil
}

func (model MessageModelImpl) Count(conversationId int64) (int, error) {
	query := `SELECT COUNT(*) FROM Message WHERE conversation = ?;`
	var count int

	if err := model.Db.QueryRow(query, conversationId).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
