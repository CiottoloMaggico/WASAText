package database

import (
	"database/sql"
	"encoding/json"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
)

type Message struct {
	Id          uint64
	Conv        Conversation
	Author      User
	SendAt      string
	DeliveredAt sql.NullString
	SeenAt      sql.NullString
	ReplyTo     *uint64
	Attachment  *Image
	Content     sql.NullString

	db *appdbimpl
}

func (m Message) Status() string {
	if m.SeenAt.Valid {
		return "Seen"
	} else if m.DeliveredAt.Valid {
		return "Delivered"
	}
	return "Sent"
}

func (m Message) Validate() error {
	if err := m.Conv.Validate(); err != nil {
		return err
	}
	if err := m.Author.Validate(); err != nil {
		return err
	}
	if m.Content.Valid {
		if ok, err := validators.MessageContentValidator(m.Content.String); !ok {
			return err
		}
	}
	if m.Attachment != nil {
		if err := m.Attachment.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	var content *string = &m.Content.String
	if !m.Content.Valid {
		content = nil
	}
	return json.Marshal(&SerializedMessage{
		m.Id,
		m.Conv.GetId(),
		m.Author,
		m.Status(),
		m.ReplyTo,
		m.Attachment,
		content,
	})
}

func (db *appdbimpl) NewMessage(conv Conversation, author User, replyTo *uint64, attachment *Image, content *string) (*Message, error) {
	if err := conv.Validate(); err != nil {
		return nil, err
	}
	if err := author.Validate(); err != nil {
		return nil, err
	}

	finalContent := sql.NullString{Valid: false}
	if content != nil {
		if ok, err := validators.MessageContentValidator(*content); !ok {
			return nil, err
		}
		finalContent.Valid = true
		finalContent.String = *content
	}
	if attachment != nil {
		if err := attachment.Validate(); err != nil {
			return nil, err
		}
	}

	var messageId uint64
	var sendAt string
	if err := db.c.QueryRow(
		qCreateMessage,
		conv.GetId(),
		author.Uuid,
		replyTo,
		content,
		attachment.Uuid,
	).Scan(&messageId, &sendAt); err != nil {
		return nil, err
	}

	return &Message{
		messageId,
		conv,
		author,
		sendAt,
		sql.NullString{Valid: false},
		sql.NullString{Valid: false},
		replyTo,
		attachment,
		finalContent,
		db,
	}, nil
}

func (m Message) Delete() error {
	if err := m.Validate(); err != nil {
		return err
	}

	if _, err := m.db.c.Exec(qDeleteMessage, m.Id); err != nil {
		return err
	}
	return nil
}
