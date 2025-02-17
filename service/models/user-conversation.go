package models

import (
	"fmt"
	"github.com/ciottolomaggico/wasatext/service/database"
)

const qUserConversations = `
	SELECT
		id userConversation_id,
    	vu.user_username userConversation_name,
		'chat' userConversation_type,
		vu.image_uuid,
		vu.image_extension,
		vu.image_width,
		vu.image_height,
		vu.image_fullUrl,
		vu.image_uploadedAt
    FROM
		Chat c,
        User_Conversation uc,
        ViewUsers vu
	WHERE uc.conversation = c.id
        AND uc.user = ?
        AND (vu.user_uuid != uc.user AND (vu.user_uuid = c.user1 OR vu.user_uuid = c.user2))
UNION
	SELECT
		gc.id userConversation_id,
        gc.name userConversation_name,
        'group' userConversation_type,
        i.uuid       image_uuid,
        i.extension  image_extension,
        i.width      image_width,
        i.height     image_height,
        i.fullUrl    image_fullUrl,
        i.uploadedAt image_uploadedAt
    FROM
		User_Conversation uc,
        GroupConversation gc,
        Image i
    WHERE uc.conversation = gc.id
        AND uc.user = ?
        AND i.uuid = gc.photo
`

const qGetUserConversation = `
WITH UserConversations AS (
	` + qUserConversations + `
)
SELECT
    uc.*,
    vm.message_id,
    vm.message_sendAt message_sendAt,
	vm.message_deliveredAt,
	vm.message_seenAt,
	vm.message_replyTo,
	vm.message_content,
	vm.attachment_uuid message_attachment,
    vm.user_uuid,
	vm.user_username,
	vm.image_uuid user_image,
    CASE
        WHEN (um.status = 3) THEN true
        ELSE false
    END message_status
FROM UserConversations uc
LEFT OUTER JOIN ViewLatestMessages vm ON vm.message_conversation = uc.userConversation_id
LEFT OUTER JOIN User_Message um ON vm.message_id = um.message AND um.user = ?
WHERE uc.userConversation_id = ?;
`

const qGetUserConversations = `
WITH UserConversations AS (
	` + qUserConversations + `
)
SELECT
    uc.*,
    vm.message_id,
    vm.message_sendAt message_sendAt,
	vm.message_deliveredAt,
	vm.message_seenAt,
	vm.message_replyTo,
	vm.message_content,
	vm.attachment_uuid message_attachment,
    vm.user_uuid,
	vm.user_username,
	vm.image_uuid user_image,
    CASE
        WHEN (um.status = 3) THEN true
        ELSE false
    END message_status
FROM UserConversations uc
LEFT OUTER JOIN ViewLatestMessages vm ON vm.message_conversation = uc.userConversation_id
LEFT OUTER JOIN User_Message um ON vm.message_id = um.message AND um.user = ?
`

type UserConversation struct {
	Id   int64  `db:"userConversation_id"`
	Name string `db:"userConversation_name"`
	Type string `db:"userConversation_type"`
	Read bool   `db:"message_status"`
	Image
	UserConversationMessagePreview
}

type UserConversationModel interface {
	GetUserConversation(userUUID string, conversationId int64) (*UserConversation, error)
	GetUserConversations(userUUID string, parameters database.QueryParameters) ([]UserConversation, int64, error)
	Count(userUUID string, parameters database.QueryParameters) (int, int64, error)
}

type UserConversationModelImpl struct {
	Db database.AppDatabase
}

func (model UserConversationModelImpl) getCursor() (int64, error) {
	var cursor int64
	query := `
		SELECT ifnull(MAX(id), 0) FROM Conversation;
	`

	if err := model.Db.QueryRow(query).Scan(&cursor); err != nil {
		return -1, err
	}
	return cursor, nil
}

func (model UserConversationModelImpl) GetUserConversation(userUUID string, conversationId int64) (*UserConversation, error) {
	query := qGetUserConversation

	userConversation := UserConversation{}
	if err := model.Db.QueryStructRow(&userConversation, query, userUUID, userUUID, userUUID, conversationId); err != nil {
		return nil, err
	}
	return &userConversation, nil
}

func (model UserConversationModelImpl) GetUserConversations(userUUID string, parameters database.QueryParameters) ([]UserConversation, int64, error) {
	cursor := parameters.Cursor
	if cursor == -1 {
		tmpCursor, err := model.getCursor()
		if err != nil {
			return nil, -1, err
		}
		cursor = tmpCursor
	}

	query := fmt.Sprintf("%s WHERE userConversation_id <= %d", qGetUserConversations, cursor)
	if filter := parameters.Filter; filter != "" {
		query += fmt.Sprintf(" AND (%s)", filter)
	}
	query += " ORDER BY message_sendAt DESC, uc.userConversation_id DESC LIMIT ? OFFSET ?;"

	userConversations := make([]UserConversation, 0, parameters.Limit)
	if err := model.Db.QueryStruct(
		&userConversations, query,
		userUUID, userUUID, userUUID,
		parameters.Limit, parameters.Offset,
	); err != nil {
		return nil, -1, err
	}

	return userConversations, cursor, nil
}

func (model UserConversationModelImpl) Count(userUUID string, parameters database.QueryParameters) (int, int64, error) {
	var count int
	cursor := parameters.Cursor
	if cursor == -1 {
		tmpCursor, err := model.getCursor()
		if err != nil {
			return 0, -1, err
		}
		cursor = tmpCursor
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM (%s) WHERE userConversation_id <= %d", qGetUserConversations, cursor)
	if filter := parameters.Filter; filter != "" {
		query += fmt.Sprintf(" AND (%s)", filter)
	}
	query += ";"

	if err := model.Db.QueryRow(query, userUUID, userUUID, userUUID).Scan(&count); err != nil {
		return 0, -1, err
	}
	return count, cursor, nil
}
