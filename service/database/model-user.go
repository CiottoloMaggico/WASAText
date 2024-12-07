package database

import (
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
	Uuid         string `json:"uuid"`
	Username     string `json:"username"`
	ProfileImage Image  `json:"photo"`

	db *appdbimpl
}

type UpdateUserParams struct {
	Username     string
	ProfileImage Image
}

func (up UpdateUserParams) Validate() error {
	if ok, err := validators.UsernameIsValid(up.Username); !ok {
		return err
	}
	if err := up.ProfileImage.Validate(); err != nil {
		return err
	}
	return nil
}

func (u User) Validate() error {
	if ok, err := validators.UsernameIsValid(u.Username); !ok {
		return err
	}
	if err := u.ProfileImage.Validate(); err != nil {
		return err
	}
	return nil
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

func (u *User) Update(params UpdateUserParams) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if err := params.Validate(); err != nil {
		return err
	}

	if _, err := u.db.c.Exec(
		qUpdateUser,
		params.Username,
		params.ProfileImage.Uuid,
		u.Uuid,
	); err != nil {
		return err
	}

	u.Username, u.ProfileImage = params.Username, params.ProfileImage
	return nil
}

func (u *User) GetConversations(pageSize int, pageNumber int) ([]DefaultConversation, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}
	rows, err := u.db.c.Query(
		getUserConversationsPaginated,
		u.Uuid, pageSize, pageNumber*pageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]DefaultConversation, 0, pageSize)
	for rows.Next() {
		var chat Conversation
		var groupName, groupPhoto, user1Id, user2Id sql.NullString
		var convId uint64

		if err := rows.Scan(
			&convId,
			&user1Id,
			&user2Id,
			&groupName,
			&groupPhoto,
		); err != nil {
			return nil, err
		}

		if groupName.Valid {
			photo, _ := u.db.GetImage(groupPhoto.String)
			chat = GroupConversation{convId, groupName.String, *photo, *u, u.db}
		} else {
			user1, _ := u.db.GetUserByUUID(user1Id.String)
			user2, _ := u.db.GetUserByUUID(user2Id.String)
			switch u.Uuid {
			case user1.Uuid:
				chat = ChatConversation{convId, *user1, *user2, *u, u.db}
			case user2.Uuid:
				chat = ChatConversation{convId, *user2, *user1, *u, u.db}
			}
		}

		chats = append(chats, DefaultConversation{chat})
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

func (db *appdbimpl) getUser(identifier string, query string) (*User, error) {
	user := User{}
	user.ProfileImage = Image{}
	image := &user.ProfileImage
	if err := db.c.QueryRow(query, identifier).Scan(
		&user.Uuid,
		&user.Username,
		&image.Uuid,
		&image.Extension,
	); err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	user.db = db
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
		user := User{}
		user.ProfileImage = Image{}
		image := &user.ProfileImage

		if err := rows.Scan(
			&user.Uuid,
			&user.Username,
			&image.Uuid,
			&image.Extension,
		); err != nil {
			return nil, err
		}

		user.db = db
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

func (u *User) GetConversation(id int64) (Conversation, error) {
	return nil, nil
	//// TODO: handling errors
	//if id < 0 {
	//	return nil, errors.New("invalid id")
	//}
	//
	//var chat Conversation
	//var groupName, groupPhoto, user1Id, user2Id sql.NullString
	//var convId int64
	//
	//if err := u.db.c.QueryRow(qGetConversation, id, u.Uuid).Scan(
	//	&convId,
	//	&user1Id,
	//	&user2Id,
	//	&groupName,
	//	&groupPhoto,
	//); errors.Is(err, sql.ErrNoRows) {
	//	return nil, errors.New("conversation not found")
	//} else if err != nil {
	//	return nil, err
	//}
	//
	//if groupName.Valid {
	//	photo, _ := u.db.GetImage(groupPhoto.String)
	//	chat = GroupConversation{convId, groupName.String, *photo, u.db}
	//} else {
	//	user1, _ := u.db.GetUserByUUID(user1Id.String)
	//	user2, _ := u.db.GetUserByUUID(user2Id.String)
	//	switch u.Uuid {
	//	case user1.Uuid:
	//		chat = ChatConversation{convId, *user1, *user2, u.db}
	//	case user2.Uuid:
	//		chat = ChatConversation{convId, *user2, *user1, u.db}
	//	}
	//}
	//
	//return chat, nil
}

func (u *User) SetDelivered() error {
	// TODO: implement
	return nil
}
