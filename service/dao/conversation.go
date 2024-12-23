package dao

import (
	"database/sql"
	"errors"
	"github.com/ciottolomaggico/wasatext/service/database"
)

type ConversationDao interface {
	CreateGroup(name string, author string, photo *string) (*database.Group, error)
	CreateChat(user1 string, user2 string) (*database.Chat, error)
	UpdateGroupName(id int64, newName string) (*database.Group, error)
	UpdateGroupPic(id int64, newPic string) (*database.Group, error)
	GetConversationParticipants(conversation int64) (database.UserWithImageList, error)
	IsParticipant(conversation int64, userUUID string) (bool, error)
	AddGroupParticipant(user string, conversation int64) error
	RemoveGroupParticipant(user string, conversation int64) error
}

type ConversationDaoImpl struct {
	db database.BaseDatabase
}

func (dao ConversationDaoImpl) CreateGroup(name string, author string, photo *string) (*database.Group, error) {
	var query string
	var params []any
	params = append(params, name, author)

	if photo != nil {
		query = "INSERT INTO GroupConversation (name, author, photo) VALUES (?,?,?) RETURNING *;"
		params = append(params, *photo)
	} else {
		query = "INSERT INTO GroupConversation (name, author) VALUES (?, ?) RETURNING *;"
	}

	return dao.db.QueryRowGroup(query, params...)
}

func (dao ConversationDaoImpl) CreateChat(user1 string, user2 string) (*database.Chat, error) {
	query := `INSERT INTO Chat (user1, user2) VALUES (?, ?) RETURNING *;`
	return dao.db.QueryRowChat(query, user1, user2)
}

func (dao ConversationDaoImpl) UpdateGroupName(id int64, newName string) (*database.Group, error) {
	query := `UPDATE GroupConversation SET name = ? WHERE id = ? RETURNING *;`

	return dao.db.QueryRowGroup(query, newName, id)
}

func (dao ConversationDaoImpl) UpdateGroupPic(id int64, newPic string) (*database.Group, error) {
	query := `UPDATE GroupConversation SET photo = ? WHERE id = ? RETURNING *;`

	return dao.db.QueryRowGroup(query, newPic, id)
}

func (dao ConversationDaoImpl) GetConversationParticipants(conversation int64) (database.UserWithImageList, error) {
	query := `SELECT u.* FROM ViewUsers u, User_Conversation uc WHERE uc.conversation = ? AND uc.user = u.userUUID;`

	return dao.db.QueryUserWithImage(query, conversation)
}

func (dao ConversationDaoImpl) IsParticipant(conversation int64, userUUID string) (bool, error) {
	query := `SELECT u.* FROM ViewUsers u, User_Conversation uc WHERE uc.conversation = ? AND u.userUUID = ? AND uc.user = u.userUUID;`

	if _, err := dao.db.QueryRowUserWithImage(query, conversation, userUUID); errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (dao ConversationDaoImpl) AddGroupParticipant(user string, conversation int64) error {
	_, err := dao.db.Exec(
		`INSERT INTO User_Conversation VALUES (?, ?);`,
		user, conversation,
	)
	return err
}

func (dao ConversationDaoImpl) RemoveGroupParticipant(user string, conversation int64) error {
	_, err := dao.db.Exec(
		`DELETE FROM User_Conversation WHERE user = ? AND conversation = ?;`,
		user, conversation,
	)
	return err
}
