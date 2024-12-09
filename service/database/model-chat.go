package database

import (
	"context"
	"errors"
)

var InvalidAuthorOrRecipient = errors.New("invalid author or recipient")

type Chat struct {
	id    uint64
	user1 User
	user2 User

	context context.Context
	db      *appdbimpl
}

func (c Chat) GetId() uint64 {
	return c.id
}

func (c Chat) GetName() string {
	requestIssuer := c.context.Value("authedUser").(User)
	switch requestIssuer.GetUUID() {
	case c.user1.GetUUID():
		return c.user2.GetUsername()
	case c.user2.GetUUID():
		return c.user1.GetUsername()
	}
	return ""
}

func (c Chat) GetPhoto() Image {
	requestIssuer := c.context.Value("authedUser").(User)
	switch requestIssuer.GetUUID() {
	case c.user1.GetUUID():
		return c.user2.GetImage()
	case c.user2.GetUUID():
		return c.user1.GetImage()
	}
	return Image{"default_user_image", ".jpg", c.db}
}

func (c Chat) GetType() string {
	return "one_to_one"
}

func (c Chat) GetContext() context.Context {
	return c.context
}

func (c Chat) GetDB() *appdbimpl {
	return c.db
}

func (c Chat) GetParticipants() ([]User, error) {
	return []User{c.user1, c.user2}, nil
}

func (c Chat) Validate() error {
	if err := c.user1.Validate(); err != nil {
		return err
	}
	if err := c.user2.Validate(); err != nil {
		return err
	}
	return nil
}

// Chat conversation create
func (db *appdbimpl) NewChat(recipient User, ctx context.Context) (*Chat, error) {
	author := ctx.Value("authedUser").(User)

	tx, err := db.c.Begin()
	if err != nil {
		return nil, err
	}

	var conversationId uint64
	if err := tx.QueryRow(qCreateConversation).Scan(&conversationId); err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := tx.Exec(qCreateChat, conversationId, author.GetUUID(), recipient.GetUUID()); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	chat := Chat{
		conversationId, author, recipient, ctx, db,
	}

	return &chat, nil
}

func (db *appdbimpl) GetChat(conversationId uint64, ctx context.Context) (BaseConversation, error) {
	chat := Chat{
		user1:   User{image: Image{}, db: db},
		user2:   User{image: Image{}, db: db},
		context: ctx,
		db:      db,
	}

	if err := db.c.QueryRow(qGetChat, conversationId).Scan(
		&chat.id,
		&chat.user1.uuid,
		&chat.user1.username,
		&chat.user1.image.uuid,
		&chat.user1.image.extension,
		&chat.user2.uuid,
		&chat.user2.username,
		&chat.user2.image.uuid,
		&chat.user2.image.extension,
	); err != nil {
		return nil, err
	}

	return chat, nil

}
