package database

import (
	"context"
	"database/sql"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
)

type Message struct {
	id           uint64
	conversation uint64
	author       User
	sendAt       string
	deliveredAt  sql.NullString
	seenAt       sql.NullString
	replyTo      *uint64
	attachment   *Image
	content      sql.NullString

	context context.Context
	db      *appdbimpl
}

func (m Message) GetId() uint64 {
	return m.id
}

func (m Message) GetConversation() (BaseConversation, error) {
	authedUser := m.context.Value("authedUser").(User)
	return authedUser.GetConversation(m.conversation, m.context)
}

func (m Message) GetAuthor() User {
	return m.author
}

func (m Message) GetTimestamps() (sendAt string, deliveredAt *string, seenAt *string) {
	sendAt = m.sendAt

	if m.deliveredAt.Valid {
		deliveredAt = &m.deliveredAt.String
	}
	if m.seenAt.Valid {
		seenAt = &m.seenAt.String
	}

	return
}

func (m Message) GetRepliedMessage() (*Message, error) {
	if m.replyTo == nil {
		return nil, nil
	}
	return m.db.GetMessage(*m.replyTo, m.context)
}

func (m Message) GetAttachment() *Image {
	return m.attachment
}

func (m Message) GetContent() *string {
	if m.content.Valid {
		return &m.content.String
	}
	return nil
}

func (m Message) ReadByUser() (bool, error) {
	requestIssuer := m.context.Value("authedUser").(User)
	var read bool
	if err := m.db.c.QueryRow(qMessageIsRead, requestIssuer.GetUUID(), m.id).Scan(&read); err != nil {
		return false, err
	}
	return read, nil
}

func (m Message) Delete() error {
	if _, err := m.db.c.Exec(qDeleteMessage, m.id); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) NewMessage(conv Conversation, replyTo *uint64, attachment *Image, content *string, ctx context.Context) (*Message, error) {
	if err := conv.Validate(); err != nil {
		return nil, err
	}
	if attachment != nil {
		if err := attachment.Validate(); err != nil {
			return nil, err
		}
	}

	messageContent := sql.NullString{Valid: false}
	if content != nil {
		if ok, err := validators.MessageContentValidator(*content); !ok {
			return nil, err
		}
		messageContent.Valid, messageContent.String = true, *content
	}
	author := ctx.Value("authedUser").(User)

	var messageId uint64
	var sendAt string
	if err := db.c.QueryRow(
		qCreateMessage,
		conv.GetId(),
		author.GetUUID(),
		replyTo,
		content,
		attachment.GetUUID(),
	).Scan(&messageId, &sendAt); err != nil {
		return nil, err
	}

	return &Message{
		messageId,
		conv.GetId(),
		author,
		sendAt,
		sql.NullString{Valid: false},
		sql.NullString{Valid: false},
		replyTo,
		attachment,
		messageContent,
		ctx,
		db,
	}, nil
}

func (db *appdbimpl) GetMessage(messageId uint64, ctx context.Context) (*Message, error) {
	message := Message{author: User{image: Image{}}, context: ctx, db: db}
	author := &message.author
	authorImage := &message.author.image
	var attachmentUUID, attachmentExt sql.NullString
	if err := db.c.QueryRow(
		qGetMessage, messageId,
	).Scan(
		&message.id,
		&message.conversation,
		&message.sendAt,
		&message.deliveredAt,
		&message.seenAt,
		&message.replyTo,
		&attachmentUUID,
		&attachmentExt,
		&author.uuid,
		&author.username,
		&authorImage.uuid,
		&authorImage.extension,
	); err != nil {
		return nil, err
	}

	if attachmentUUID.Valid {
		message.attachment = &Image{attachmentUUID.String, attachmentExt.String, db}
	}

	return &message, nil
}
