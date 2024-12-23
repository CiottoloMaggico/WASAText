package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ciottolomaggico/wasatext/service/utils/authentication"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
	"github.com/google/uuid"
	"net/http"
)

const defaultUserImage = "default_user_image"

type User struct {
	uuid     string
	username string
	image    Image

	db *appdbimpl
}

func (u User) GetUUID() string {
	return u.uuid
}

func (u User) GetUsername() string {
	return u.username
}

func (u User) GetImage() Image {
	return u.image
}

func (u User) SetUsername(username string) error {
	if ok, err := validators.UsernameIsValid(username); !ok {
		return err
	}

	if _, err := u.db.c.Exec(qSetUsername, username, u.uuid); err != nil {
		return err
	}

	u.username = username
	return nil
}

func (u User) SetImage(image Image) error {
	if err := image.Validate(); err != nil {
		return err
	}

	if _, err := u.db.c.Exec(qSetPhoto, image.GetUUID(), u.uuid); err != nil {
		return err
	}

	u.image = image
	return nil
}

func (u User) Validate() error {
	if ok, err := validators.UsernameIsValid(u.username); !ok {
		return err
	}
	if err := u.image.Validate(); err != nil {
		return err
	}
	return nil
}

func (u User) GetConversation(conversationId uint64, ctx context.Context) (BaseConversation, error) {
	var user1UUID, user2UUID, groupName, groupPhotoUUID sql.NullString
	var tmpConvId uint64

	if err := u.db.c.QueryRow(qGetConversation, conversationId, u.uuid).Scan(
		&tmpConvId,
		&user1UUID,
		&user2UUID,
		&groupName,
		&groupPhotoUUID,
	); err != nil {
		return nil, err
	}

	if user1UUID.Valid && user2UUID.Valid {
		return u.db.GetChat(conversationId, ctx)
	} else {
		return u.db.GetGroup(conversationId, ctx)
	}
	return nil, nil
}

func (u User) GetConversations(pageSize int, pageNumber int, ctx context.Context) ([]BaseConversation, error) {
	rows, err := u.db.c.Query(
		getUserConversationsPaginated,
		u.uuid, pageSize, pageNumber*pageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]BaseConversation, 0, pageSize)
	for rows.Next() {
		var convId uint64
		var user1UUID, user2UUID, groupName, groupPhotoUUID sql.NullString
		var chat BaseConversation

		if err := rows.Scan(
			&convId,
			&user1UUID,
			&user2UUID,
			&groupName,
			&groupPhotoUUID,
		); err != nil {
			return nil, err
		}

		if user1UUID.Valid && user2UUID.Valid {
			chat, err = u.db.GetChat(convId, ctx)
			if err != nil {
				return nil, err
			}
		} else {
			chat, err = u.db.GetGroup(convId, ctx)
			if err != nil {
				return nil, err
			}
		}

		chats = append(chats, chat)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); errors.Is(err, sql.ErrNoRows) {
		return chats, nil
	} else if err != nil {
		return nil, err
	}

	return chats, nil
}

func (db *appdbimpl) NewUser(username string) (*User, error) {
	// For each provided arguments run the corresponding validator
	if ok, err := validators.UsernameIsValid(username); !ok {
		return nil, err
	}

	// If all the arguments are valid then set the "private" object fields (e.g. primary key)
	rawUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate a new uuid: %w", err)
	}

	// Create the user in the database
	if _, err := db.c.Exec(
		qCreateUser, rawUUID.String(), username, defaultUserImage,
	); err != nil {
		return nil, fmt.Errorf("failed to save the user in the db: %w", err)
	}

	user := User{rawUUID.String(), username, Image{defaultUserImage, ".jpg", db}, db}
	return &user, nil
}

func (db *appdbimpl) getUser(identifier string, query string) (*User, error) {
	user := User{image: Image{}, db: db}
	image := &user.image
	if err := db.c.QueryRow(query, identifier).Scan(
		&user.uuid,
		&user.username,
		&image.uuid,
		&image.extension,
	); err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

func (db *appdbimpl) GetUserByUsername(username string) (*User, error) {
	if ok, err := validators.UsernameIsValid(username); !ok {
		return nil, err
	}
	return db.getUser(username, qGetUserByUsername)
}

func (db *appdbimpl) GetUserByUUID(UUID string) (*User, error) {
	if _, err := uuid.Parse(UUID); err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}
	return db.getUser(UUID, qGetUserByUUID)
}

// TODO: maybe in controller?
func (db *appdbimpl) GetAuthenticatedUser(r *http.Request) (*User, error) {
	token := authentication.GetAuthToken(r)
	if _, err := uuid.Parse(token); err != nil {
		return nil, errors.New("Invalid token or User not authenticated")
	}

	return db.GetUserByUUID(token)
}

func (db *appdbimpl) GetUsers(pageSize int, pageNumber int) ([]User, error) {
	rows, err := db.c.Query(
		qGetUsersPaginated,
		pageSize, pageNumber*pageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0, pageSize)
	for rows.Next() {
		user := User{image: Image{}, db: db}
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

func (db *appdbimpl) UsersCount() (int, error) {
	var count int
	err := db.c.QueryRow(
		qGetUsersCount,
	).Scan(
		&count,
	)
	return count, err
}

func (u *User) SetDelivered() error {
	// TODO: implement
	return nil
}
