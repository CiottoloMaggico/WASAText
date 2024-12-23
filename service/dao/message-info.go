package dao

import (
	"github.com/ciottolomaggico/wasatext/service/database"
)

type MessageInfoDao interface {
	GetMessageComments(messageId int64) (database.MessageInfoList, error)
	GetMessageInfo(user string, message int64) (*database.MessageInfo, error)
	SetComment(user string, message int64, comment string) (*database.MessageInfo, error)
	RemoveComment(user string, message int64) error
}

type MessageInfoDaoImpl struct {
	db database.BaseDatabase
}

func (dao MessageInfoDaoImpl) GetMessageComments(messageId int64) (database.MessageInfoList, error) {
	query := `SELECT * FROM User_Message WHERE message = ? AND comment IS NOT NULL`

	return dao.db.QueryMessageInfo(query, messageId)
}

func (dao MessageInfoDaoImpl) GetMessageInfo(user string, message int64) (*database.MessageInfo, error) {
	query := `SELECT * FROM User_Message WHERE user = ? AND message = ?`
	return dao.db.QueryRowMessageInfo(query, user, message)
}

func (dao MessageInfoDaoImpl) SetComment(user string, message int64, comment string) (*database.MessageInfo, error) {
	query := `UPDATE User_Message SET comment = ? WHERE user = ? AND message = ? RETURNING *;`
	return dao.db.QueryRowMessageInfo(query, comment, user, message)
}

func (dao MessageInfoDaoImpl) RemoveComment(user string, message int64) error {
	query := `UPDATE User_Message SET comment = NULL WHERE user = ? AND message = ? RETURNING *;`
	_, err := dao.db.QueryRowMessageInfo(query, user, message)
	return err
}
