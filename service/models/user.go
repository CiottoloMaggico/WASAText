package models

import (
	"fmt"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/gofrs/uuid"
)

type UserModel interface {
	CreateUser(username string) (*User, error)
	UpdateUsername(uuid string, newUsername string) (*User, error)
	UpdateProfilePic(uuid string, newPhoto string) (*User, error)
	GetUserWithImage(uuid string) (*UserWithImage, error)
	GetUserWithImageByUsername(username string) (*UserWithImage, error)
	GetUsersWithImage(parameters database.QueryParameters) ([]UserWithImage, error)
	Count(parameters database.QueryParameters) (int, error)
}

type User struct {
	Uuid     string `db:"user_uuid"`
	Username string `db:"user_username"`
	Photo    string `db:"user_photo"`
}

type UserWithImage struct {
	Uuid     string `db:"user_uuid"`
	Username string `db:"user_username"`
	Image
}

type UserModelImpl struct {
	Db database.AppDatabase
}

func (model UserModelImpl) Count(parameters database.QueryParameters) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM ViewUsers`

	if filter := parameters.Filter; filter != "" {
		query += fmt.Sprintf(" WHERE (%s)", filter)
	}

	query += ";"
	if err := model.Db.QueryRow(query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (model UserModelImpl) CreateUser(username string) (*User, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO User (uuid, username)
		VALUES (?, ?)
		RETURNING uuid user_uuid, username user_username, photo user_photo;
	`

	user := User{}
	if err := model.Db.QueryStructRow(&user, query, newUUID.String(), username); err != nil {
		return nil, err
	}

	return &user, nil
}

func (model UserModelImpl) UpdateUsername(uuid string, newUsername string) (*User, error) {
	query := `
		UPDATE User
		SET username = ?
		WHERE uuid = ?
		RETURNING uuid user_uuid, username user_username, photo user_photo;
	`

	user := User{}
	if err := model.Db.QueryStructRow(&user, query, newUsername, uuid); err != nil {
		return nil, err
	}

	return &user, nil
}

func (model UserModelImpl) UpdateProfilePic(uuid string, newPhoto string) (*User, error) {
	query := `
		UPDATE User
		SET photo = ?
		WHERE uuid = ?
		RETURNING uuid user_uuid, username user_username, photo user_photo;
	`

	user := User{}
	if err := model.Db.QueryStructRow(&user, query, newPhoto, uuid); err != nil {
		return nil, err
	}

	return &user, nil
}

func (model UserModelImpl) GetUserWithImage(uuid string) (*UserWithImage, error) {
	query := `SELECT * FROM ViewUsers WHERE user_uuid = ?;`

	user := UserWithImage{}
	if err := model.Db.QueryStructRow(&user, query, uuid); err != nil {
		return nil, err
	}

	return &user, nil
}

func (model UserModelImpl) GetUserWithImageByUsername(username string) (*UserWithImage, error) {
	query := `SELECT * FROM ViewUsers WHERE user_username = ?;`

	user := UserWithImage{}
	if err := model.Db.QueryStructRow(&user, query, username); err != nil {
		return nil, err
	}

	return &user, nil
}

func (model UserModelImpl) GetUsersWithImage(parameters database.QueryParameters) ([]UserWithImage, error) {
	query := `SELECT * FROM ViewUsers`

	if filter := parameters.Filter; filter != "" {
		query += fmt.Sprintf(" WHERE (%s)", filter)
	}

	query += " LIMIT ? OFFSET ?;"

	users := make([]UserWithImage, 0, parameters.Limit)
	if err := model.Db.QueryStruct(&users, query, parameters.Limit, parameters.Offset); err != nil {
		return nil, err
	}

	return users, nil
}
