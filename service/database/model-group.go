package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
)

const MAX_GROUP_SIZE = 200

var InvalidId = errors.New("invalid id")
var NotMember = errors.New("this user isn't a group member")
var ExceedMaxGroupSize = fmt.Errorf("groups can handle at most %i participants", MAX_GROUP_SIZE)

// TODO: update request issuer with context
type Group struct {
	id     uint64
	name   string
	photo  Image
	author User

	context context.Context
	db      *appdbimpl
}

func (g Group) GetId() uint64 {
	return g.id
}

func (g Group) GetName() string {
	return g.name
}

func (g Group) SetName(name string) error {
	if ok, err := validators.GroupNameIsValid(name); !ok {
		return err
	}

	if _, err := g.db.c.Exec(qSetGroupName, name, g.id); err != nil {
		return err
	}

	g.name = name
	return nil
}

func (g Group) GetPhoto() Image {
	return g.photo
}

func (g Group) SetPhoto(photo Image) error {
	if _, err := g.db.c.Exec(qSetGroupPhoto, photo.GetUUID(), g.id); err != nil {
		return err
	}

	g.photo = photo
	return nil
}

func (g Group) GetType() string {
	return "group"
}

func (g Group) GetAuthor() User {
	return g.author
}

func (g Group) GetContext() context.Context {
	return g.context
}

func (g Group) GetDB() *appdbimpl {
	return g.db
}

func (g Group) Validate() error {
	if ok, err := validators.GroupNameIsValid(g.name); !ok {
		return err
	}
	if err := g.photo.Validate(); err != nil {
		return err
	}
	if err := g.author.Validate(); err != nil {
		return err
	}

	return nil
}

func (g Group) GetParticipants() ([]User, error) {
	rows, err := g.db.c.Query(
		qGetConversationParticipants, g.GetId(),
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0, MAX_GROUP_SIZE)
	for rows.Next() {
		user := User{image: Image{}, db: g.db}
		image := &user.image

		if err := rows.Scan(
			&user.uuid,
			&user.username,
			&image.uuid,
			&image.extension,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (g Group) AddParticipant(user User) error {
	if _, err := g.db.c.Exec(qLinkUserToConversation, user.GetUUID(), g.id); err != nil {
		return err
	}
	return nil
}

func (g Group) RemoveParticipant(user User) error {
	if _, err := g.db.c.Exec(qUnLinkUserToConversation, user.GetUUID(), g.id); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) NewGroup(name string, photo Image, ctx context.Context) (*Group, error) {
	// TODO: remove the feature to add participants in the same request of group creation, rather handle it with different requests in frontend
	if ok, err := validators.GroupNameIsValid(name); !ok {
		return nil, err
	}
	if err := photo.Validate(); err != nil {
		return nil, err
	}

	author := ctx.Value("authedUser").(User)

	tx, err := db.c.Begin()
	if err != nil {
		return nil, err
	}

	// Create the conversation
	var conversationId uint64
	if err := tx.QueryRow(qCreateConversation).Scan(&conversationId); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create a new group with the same id of the conversation
	if _, err := tx.Exec(qCreateGroup, conversationId, name, author.GetUUID(), photo.GetUUID()); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	group := Group{
		conversationId, name, photo, author, ctx, db,
	}

	return &group, nil
}

func (db *appdbimpl) GetGroup(conversationId uint64, ctx context.Context) (GroupConversation, error) {
	group := Group{
		photo:   Image{db: db},
		author:  User{image: Image{}, db: db},
		context: ctx,
		db:      db,
	}

	if err := db.c.QueryRow(qGetGroup, conversationId).Scan(
		&group.id,
		&group.name,
		&group.author.uuid,
		&group.author.username,
		&group.author.image.uuid,
		&group.author.image.extension,
		&group.photo.uuid,
		&group.photo.extension,
	); err != nil {
		return nil, err
	}

	return group, nil
}
