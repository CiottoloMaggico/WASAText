package database

import (
	"context"
	"database/sql"
	"errors"
)

type BaseConversation interface {
	GetId() uint64
	GetName() string
	GetPhoto() Image
	GetType() string
	GetContext() context.Context
	GetDB() *appdbimpl
	GetParticipants() ([]User, error)
	Validate() error
}

type GroupConversation interface {
	BaseConversation
	GetAuthor() User
	AddParticipant(user User) error
	RemoveParticipant(user User) error
}

type Conversation struct {
	BaseConversation
}

func (c Conversation) SendMessage(replyTo *uint64, attachment *Image, content *string) (*Message, error) {
	return c.GetDB().NewMessage(c, replyTo, attachment, content, c.GetContext())
}

func (c Conversation) GetMessage(messageId uint64) (*Message, error) {
	return c.GetDB().GetMessage(messageId, c.GetContext())
}

func (dc *Conversation) GetMessages(pageSize int, pageNumber int) ([]Message, error) {
	conversationId, db, ctx := dc.GetId(), dc.GetDB(), dc.GetContext()
	rows, err := db.c.Query(
		qGetConversationMessagesPaginated,
		dc.GetId(), pageSize, pageNumber*pageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]Message, 0, pageSize)
	for rows.Next() {
		message := Message{
			conversation: conversationId,
			author:       User{image: Image{}},
			context:      ctx,
			db:           db,
		}
		author := &message.author
		authorImage := &message.author.image
		var attachmentUUID, attachmentExt sql.NullString
		if err := rows.Scan(
			&message.id,
			&message.sendAt,
			&message.deliveredAt,
			&message.seenAt,
			&message.replyTo,
			&message.content,
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

		messages = append(messages, message)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); errors.Is(err, sql.ErrNoRows) {
		return messages, nil
	} else if err != nil {
		return nil, err
	}

	return messages, nil
}

func (dc *Conversation) GetLatestMessage() (*Message, error) {
	conversationId, db, ctx := dc.GetId(), dc.GetDB(), dc.GetContext()
	message := Message{
		conversation: conversationId,
		author:       User{image: Image{}},
		context:      ctx,
		db:           db,
	}
	author := &message.author
	authorImage := &message.author.image
	var attachmentUUID, attachmentExt sql.NullString
	if err := db.c.QueryRow(qGetLatestMessage, conversationId).Scan(
		&message.id,
		&message.sendAt,
		&message.deliveredAt,
		&message.seenAt,
		&message.replyTo,
		&message.content,
		&attachmentUUID,
		&attachmentExt,
		&author.uuid,
		&author.username,
		&authorImage.uuid,
		&authorImage.extension,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if attachmentUUID.Valid {
		message.attachment = &Image{attachmentUUID.String, attachmentExt.String, db}
	}
	return &message, nil
}

// set message as seen
