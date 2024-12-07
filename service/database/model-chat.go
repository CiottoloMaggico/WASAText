package database

import (
	"errors"
)

var InvalidAuthorOrRecipient = errors.New("invalid author or recipient")

type ChatConversation struct {
	Id            uint64
	RequestSender User
	Recipient     User

	requestIssuer User
	db            *appdbimpl
}

func (c ChatConversation) GetId() uint64 {
	return c.Id
}

func (c ChatConversation) GetName() string {
	return c.Recipient.Username
}

func (c ChatConversation) GetPhoto() Image {
	return c.Recipient.ProfileImage
}

func (c ChatConversation) GetRequestIssuer() User {
	return c.requestIssuer
}

func (c ChatConversation) GetDB() *appdbimpl {
	return c.db
}

func (c ChatConversation) Type() string {
	return "one_to_one"
}

func (c ChatConversation) Validate() error {
	if c.Id < 0 {
		return InvalidId
	}
	if err := c.RequestSender.Validate(); err != nil {
		return err
	}
	if err := c.Recipient.Validate(); err != nil {
		return err
	}
	return nil
}

// Chat conversation create
func (db *appdbimpl) NewChat(author *User, recipient *User) (*ChatConversation, error) {
	if author.Validate() != nil || recipient.Validate() != nil {
		return nil, InvalidAuthorOrRecipient
	}

	tx, err := db.c.Begin()
	if err != nil {
		return nil, err
	}

	var conversationId uint64
	if err := tx.QueryRow(qCreateConversation).Scan(&conversationId); err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := tx.Exec(qCreateChat, conversationId, author.Uuid, recipient.Uuid); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &ChatConversation{
		conversationId, *author, *recipient, *author, db,
	}, nil
}
