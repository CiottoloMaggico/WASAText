package dao

import (
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/gofrs/uuid"
)

type UserDao interface {
	CreateUser(username string) (*database.User, error)
	UpdateUsername(uuid string, newUsername string) (*database.User, error)
	UpdateProfilePic(uuid string, newPhoto string) (*database.User, error)
	GetUser(uuid string) (*database.UserWithImage, error)
	GetUserByUsername(username string) (*database.UserWithImage, error)
	GetUsers(page uint, size uint) (database.UserWithImageList, error)
}

type UserDaoImpl struct {
	db database.BaseDatabase
}

func (dao UserDaoImpl) CreateUser(username string) (*database.User, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO User (uuid, username)  VALUES (?, ?) RETURNING *;`

	return dao.db.QueryRowUser(query, newUUID.String(), username)
}

func (dao UserDaoImpl) UpdateUsername(uuid string, newUsername string) (*database.User, error) {
	query := `UPDATE User SET username = ? WHERE uuid = ? RETURNING *;`

	return dao.db.QueryRowUser(query, newUsername, uuid)
}

func (dao UserDaoImpl) UpdateProfilePic(uuid string, newPhoto string) (*database.User, error) {
	query := `UPDATE User SET photo = ? WHERE uuid = ? RETURNING *;`

	return dao.db.QueryRowUser(query, newPhoto, uuid)
}

func (dao UserDaoImpl) GetUser(uuid string) (*database.UserWithImage, error) {
	query := `SELECT * FROM ViewUsers WHERE userUUID = ?;`

	return dao.db.QueryRowUserWithImage(query, uuid)
}

func (dao UserDaoImpl) GetUserByUsername(username string) (*database.UserWithImage, error) {
	query := `SELECT * FROM ViewUsers WHERE username = ?;`

	return dao.db.QueryRowUserWithImage(query, username)
}

func (dao UserDaoImpl) GetUsers(page uint, size uint) (database.UserWithImageList, error) {
	query := `SELECT * FROM ViewUsers ORDER BY rowid LIMIT ? OFFSET ?;`

	return dao.db.QueryUserWithImage(query, size, page*size)
}
