package models

import (
	"github.com/ciottolomaggico/wasatext/service/database"
)

const qUserConversations = `
	SELECT
		id userConversation_id,
    	vu.username             userConversation_name,
		'chat' userConversation_type,
		vu.iUuid       image_uuid,
		vu.extension        image_extension,
		vu.profilePicWidth      image_width,
		vu.profilePicHeight     image_height,
		vu.profilePicFullUrl    image_fullUrl,
		vu.profilePicUploadedAt image_uploadedAt
    FROM
		Chat c,
        User_Conversation uc,
        ViewUsers vu
	WHERE uc.conversation = c.id
        AND uc.user = ?
        AND (vu.uUuid != uc.user AND (vu.uUuid = c.user1 OR vu.uUuid = c.user2))
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
    vm.id message_id,
    MAX(vm.sendAt) sendAt message_sendAt,
	vm.deliveredAt message_deliveredAt,
	vm.seenAt message_seenAt,
	vm.replyTo message_replyTo,
	vm.content message_content,
	vm.attachmentUUID message_attachment,
    vm.userUUID user_uuid,
	vm.username user_username,
	vm.profilePicUUID user_photo,
    CASE WHEN um.status = 3 THEN true ELSE false END message_status
FROM UserConversations uc
LEFT OUTER JOIN ViewMessages vm ON vm.conversation = uc.id
LEFT OUTER JOIN User_Message um ON vm.id = um.message AND um.user = ?
WHERE uc.id = ?;
`

const qGetUserConversations = `
WITH UserConversations AS (
	` + qUserConversations + `
)
SELECT
    uc.*,
    vm.id message_id,
    MAX(vm.sendAt) sendAt message_sendAt,
	vm.deliveredAt message_deliveredAt,
	vm.seenAt message_seenAt,
	vm.replyTo message_replyTo,
	vm.content message_content,
	vm.attachmentUUID message_attachment,
    vm.userUUID user_uuid,
	vm.username user_username,
	vm.profilePicUUID user_photo,
    CASE WHEN um.status = 3 THEN true ELSE false END message_status
FROM UserConversations uc
LEFT OUTER JOIN ViewMessages vm ON vm.conversation = uc.id
LEFT OUTER JOIN User_Message um ON vm.id = um.message AND um.user = ?
GROUP BY uc.id
ORDER BY um.status ASC
LIMIT ? OFFSET ?;
`

type UserConversation struct {
	Id   int64  `db:"userConversation_id"`
	Name string `db:"userConversation_name"`
	Type string `db:"userConversation_type"`
	Read bool   `db:"message_status"`
	Image
	MessageWithAuthor
}

type UserConversationModel interface {
	GetUserConversation(userUUID string, conversationId int64) (*UserConversation, error)
	GetUserConversations(userUUID string, page uint, size uint) ([]UserConversation, error)
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

func (model UserConversationModelImpl) GetUserConversations(userUUID string, page uint, size uint) ([]UserConversation, error) {
	query := qGetUserConversations

	userConversations := make([]UserConversation, size)
	if err := model.Db.QueryStruct(&userConversations, query, userUUID, userUUID, userUUID, size, page*size); err != nil {
		return nil, err
	}
	return userConversations, nil
}
