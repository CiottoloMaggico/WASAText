package models

import (
	"github.com/ciottolomaggico/wasatext/service/database"
)

type Conversation struct {
	Id int64 `db:"conversation_id"`
}

type Chat struct {
	Conversation
	User1 string `db:"chat_user1"`
	User2 string `db:"chat_user2"`
}

type Group struct {
	Conversation
	Name   string `db:"group_name"`
	Author string `db:"group_author"`
	Photo  string `db:"group_photo"`
}

type ConversationModel interface {
	CreateGroup(name string, author string, photo *string) (*Group, error)
	CreateChat(user1 string, user2 string) (*Chat, error)
	UpdateGroupName(id int64, newName string) (*Group, error)
	UpdateGroupPic(id int64, newPic string) (*Group, error)
	GetConversationParticipants(conversation int64) ([]UserWithImage, error)
	IsParticipant(conversation int64, userUUID string) (bool, error)
	AddGroupParticipant(user string, conversation int64) error
	AddGroupParticipants(users []string, conversation int64) error
	RemoveGroupParticipant(user string, conversation int64) error
}

type ConversationModelImpl struct {
	Db database.AppDatabase
}

func (model ConversationModelImpl) CreateGroup(name string, author string, photo *string) (*Group, error) {
	tx, err := model.Db.StartTx()
	if err != nil {
		return nil, err
	}

	queryConv, conv := `INSERT INTO Conversation DEFAULT VALUES RETURNING id conversation_id`, Conversation{}
	if err := tx.QueryStructRow(&conv, queryConv); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	var query string
	var params []interface{}
	params = append(params, conv.Id, name, author)

	if photo != nil {
		query = `
			INSERT INTO GroupConversation (id, name, author, photo)
			VALUES (?,?,?,?)
			RETURNING id conversation_id, name group_name, author group_author, photo group_photo;
		`
		params = append(params, *photo)
	} else {
		query = `
			INSERT INTO GroupConversation (id, name, author)
			VALUES (?, ?, ?)
			RETURNING id conversation_id, name group_name, author group_author, photo group_photo;
		`
	}

	group := Group{}
	if err := tx.QueryStructRow(&group, query, params...); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &group, nil
}

func (model ConversationModelImpl) CreateChat(user1 string, user2 string) (*Chat, error) {
	queryConv := `INSERT INTO Conversation DEFAULT VALUES RETURNING id conversation_id`
	queryChat := `
		INSERT INTO Chat (id, user1, user2)
		VALUES (?, ?, ?)
		RETURNING id conversation_id, user1 chat_user1, user2 chat_user2;
	`

	tx, err := model.Db.StartTx()
	if err != nil {
		return nil, err
	}

	conv := Conversation{}
	if err := tx.QueryStructRow(&conv, queryConv); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	chat := Chat{}
	if err := tx.QueryStructRow(&chat, queryChat, conv.Id, user1, user2); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &chat, nil
}

func (model ConversationModelImpl) UpdateGroupName(id int64, newName string) (*Group, error) {
	query := `
		UPDATE GroupConversation
		SET name = ?
		WHERE id = ?
		RETURNING id conversation_id, name group_name, author group_author, photo group_photo;
	`

	group := Group{}
	if err := model.Db.QueryStructRow(&group, query, newName, id); err != nil {
		return nil, err
	}

	return &group, nil
}

func (model ConversationModelImpl) UpdateGroupPic(id int64, newPic string) (*Group, error) {
	query := `
		UPDATE GroupConversation
		SET photo = ?
		WHERE id = ?
		RETURNING id conversation_id, name group_name, author group_author, photo group_photo;
	`

	group := Group{}
	if err := model.Db.QueryStructRow(&group, query, newPic, id); err != nil {
		return nil, err
	}

	return &group, nil
}

func (model ConversationModelImpl) GetConversationParticipants(conversation int64) ([]UserWithImage, error) {
	query := `
		SELECT u.*
		FROM ViewUsers u, User_Conversation uc
		WHERE uc.conversation = ? AND uc.user = u.user_uuid;
	`

	users := make([]UserWithImage, 0)
	if err := model.Db.QueryStruct(&users, query, conversation); err != nil {
		return nil, err
	}

	return users, nil
}

func (model ConversationModelImpl) IsParticipant(conversation int64, userUUID string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT * FROM User_Conversation WHERE user = ? AND conversation = ?);
	`

	var exists bool
	if err := model.Db.QueryRow(query, userUUID, conversation).Scan(&exists); err != nil {
		return false, database.DBError(err)
	}

	return exists, nil
}

func (model ConversationModelImpl) AddGroupParticipant(user string, conversation int64) error {
	if _, err := model.Db.Exec(
		`INSERT INTO User_Conversation VALUES (?, ?);`,
		user, conversation,
	); err != nil {
		return database.DBError(err)
	}

	return nil
}

func (model ConversationModelImpl) AddGroupParticipants(users []string, conversation int64) error {
	tx, err := model.Db.StartTx()
	if err != nil {
		return err
	}

	query := `INSERT INTO User_Conversation VALUES (?, ?);`
	preparedQuery, err := tx.Prepare(query)
	if err != nil {
		return database.DBError(err)
	}

	for _, user := range users {
		if _, err := preparedQuery.Exec(user, conversation); err != nil {
			if err := tx.Rollback(); err != nil {
				return database.DBError(err)
			}
			return database.DBError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return database.DBError(err)
	}
	return nil
}

func (model ConversationModelImpl) RemoveGroupParticipant(user string, conversation int64) error {
	if _, err := model.Db.Exec(
		`DELETE FROM User_Conversation WHERE user = ? AND conversation = ?;`,
		user, conversation,
	); err != nil {
		return database.DBError(err)
	}
	return nil
}
