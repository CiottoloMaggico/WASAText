package dao

import (
	"github.com/ciottolomaggico/wasatext/service/database"
)

const qUserConversations = `
	SELECT
		id,
    	vu.username             name,
		'chat' type,
		vu.iUuid       imageUUID,
		vu.extension        imageExt,
		vu.profilePicWidth      imageWidth,
		vu.profilePicHeight     imageHeight,
		vu.profilePicFullUrl    imageFullUrl,
		vu.profilePicUploadedAt imageUploadedAt
    FROM
		Chat c,
        User_Conversation uc,
        ViewUsers vu
	WHERE uc.conversation = c.id
        AND uc.user = ?
        AND (vu.uUuid != uc.user AND (vu.uUuid = c.user1 OR vu.uUuid = c.user2))
UNION
	SELECT
		gc.id,
        gc.name,
        'group' type,
        i.uuid       imageUUID,
        i.extension  imageExt,
        i.width      imageWidth,
        i.height     imageHeight,
        i.fullUrl    imageFullUrl,
        i.uploadedAt imageUploadedAt
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
    vm.id messageId,
    MAX(vm.sendAt) sendAt,
	vm.deliveredAt,
	vm.seenAt,
	vm.replyTo,
	vm.content,
	vm.attachmentUUID,
    vm.userUUID,
	vm.username,
	vm.profilePicUUID,
    CASE WHEN um.status = 3 THEN true ELSE false END status
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
    vm.id messageId,
    MAX(vm.sendAt) sendAt,
	vm.deliveredAt,
	vm.seenAt,
	vm.replyTo,
	vm.content,
	vm.attachmentUUID,
    vm.userUUID,
	vm.username,
	vm.profilePicUUID,
	vm.profilePicExt,
	vm.profilePicWidth,
	vm.profilePicHeight,
	vm.profilePicFullUrl,
	vm.profilePicUploadedAt,
    CASE WHEN um.status = 3 THEN true ELSE false END status
FROM UserConversations uc
LEFT OUTER JOIN ViewMessages vm ON vm.conversation = uc.id
LEFT OUTER JOIN User_Message um ON vm.id = um.message AND um.user = ?
GROUP BY uc.id
ORDER BY um.status ASC
LIMIT ? OFFSET ?;
`

type UserConversationDao interface {
	GetUserConversation(userUUID string, conversationId int64) (*database.UserConversation, error)
	GetUserConversations(userUUID string, page uint, size uint) ([]database.UserConversation, error)
}

type UserConversationDaoImpl struct {
	db database.BaseDatabase
}

func (dao UserConversationDaoImpl) GetUserConversation(userUUID string, conversationId int64) (*database.UserConversation, error) {
	query := qGetUserConversation

	return dao.db.QueryRowUserConversation(query, userUUID, userUUID, userUUID, conversationId)
}

func (dao UserConversationDaoImpl) GetUserConversations(userUUID string, page uint, size uint) ([]database.UserConversation, error) {
	query := qGetUserConversations

	return dao.db.QueryUserConversation(query, userUUID, userUUID, userUUID, page, size)
}
