package models

import (
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
    MAX(vm.message_sendAt) message_sendAt,
	vm.message_deliveredAt,
	vm.message_seenAt,
	vm.message_replyTo,
	vm.message_content,
	vm.attachment_uuid message_attachment,
    vm.user_uuid,
	vm.user_username,
	vm.image_uuid user_image,
    CASE WHEN (um.status = 3 OR vm.message_id IS NULL) THEN true ELSE false END message_status
FROM UserConversations uc
LEFT OUTER JOIN ViewMessages vm ON vm.message_conversation = uc.userConversation_id
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
    MAX(vm.message_sendAt) message_sendAt,
	vm.message_deliveredAt,
	vm.message_seenAt,
	vm.message_replyTo,
	vm.message_content,
	vm.attachment_uuid message_attachment,
    vm.user_uuid,
	vm.user_username,
	vm.image_uuid user_image,
    CASE WHEN (um.status = 3 OR vm.message_id IS NULL) THEN true ELSE false END message_status
FROM UserConversations uc
LEFT OUTER JOIN ViewMessages vm ON vm.message_conversation = uc.userConversation_id
LEFT OUTER JOIN User_Message um ON vm.message_id = um.message AND um.user = ?
GROUP BY uc.userConversation_id
ORDER BY um.status ASC
LIMIT ? OFFSET ?;
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
	GetUserConversations(userUUID string, page int, size int) ([]UserConversation, error)
	Count(userUUID string) (int, error)
}

type UserConversationModelImpl struct {
	Db database.AppDatabase
}

func (model UserConversationModelImpl) GetUserConversation(userUUID string, conversationId int64) (*UserConversation, error) {
	query := qGetUserConversation

	userConversation := UserConversation{}
	if err := model.Db.QueryStructRow(&userConversation, query, userUUID, userUUID, userUUID, conversationId); err != nil {
		return nil, err
	}
	return &userConversation, nil

}

func (model UserConversationModelImpl) GetUserConversations(userUUID string, page int, size int) ([]UserConversation, error) {
	query := qGetUserConversations

	userConversations := make([]UserConversation, 0, size)
	if err := model.Db.QueryStruct(&userConversations, query, userUUID, userUUID, userUUID, size, (page-1)*size); err != nil {
		return nil, err
	}
	return userConversations, nil
}

func (model UserConversationModelImpl) Count(userUUID string) (int, error) {
	query := `SELECT COUNT(*) FROM User_Conversation WHERE user = ?;`
	var count int

	if err := model.Db.QueryRow(query, userUUID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
