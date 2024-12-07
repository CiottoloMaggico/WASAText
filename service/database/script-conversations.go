package database

const getUserConversationsPaginated = `
	SELECT
		vc.*
	FROM ViewConversations vc, User_Conversation uc
	WHERE uc.conversation = vc.id AND uc.user = ?
	LIMIT ? OFFSET ?;
`

const qGetConversation = `
	SELECT
		vc.*
	FROM ViewConversations vc, User_Conversation uc
	WHERE uc.conversation = vc.id AND vc.id = ? AND uc.user = ?;
`

const qGetConversationParticipants = `
	SELECT
		vu.*
	FROM ViewUsers vu, User_Conversation uc
	WHERE vu.uUuid = uc.user AND uc.conversation = ?;
`

const qCreateChat = `
	INSERT INTO Chat VALUES (?, ?, ?);
`

const qAddChatToConversation = `
	INSERT INTO User_Conversation VALUES (?, ?), (?, ?);
`

const qAddGroupToConversation = `
	INSERT INTO User_Conversation VALUES (?, ?);
`

const qGetGroupId = `
	SELECT id FROM GroupConversation WHERE id = ?;
`

const qRemoveParticipant = `
	DELETE FROM User_Conversation WHERE user = ? AND conversation = ?;
`

const qDeleteConversation = `
	DELETE FROM Conversation WHERE id = ?;
`
const qCreateConversation = `
	INSERT INTO Conversation DEFAULT VALUES RETURNING id;
`

const qCreateGroup = `
	INSERT INTO GroupConversation VALUES (?, ?, ?, ?);
`

const qUpdateGroup = `
	UPDATE GroupConversation SET name = ?, photo = ? WHERE id = ?;
`

const qGetGroup = `
	SELECT id, name, author, photo, Image.extension FROM GroupConversation, Image WHERE id = ? AND photo = Image.uuid;
`

const qGetParticipant = `
	SELECT uc.user FROM User_Conversation uc WHERE uc.user = ? AND uc.conversation = ?;
`

const qGetLatestMessage = `
	SELECT
		vlm.messageId,
		vlm.sendAt,
		vlm.deliveredAt,
		vlm.seenAt,
		vlm.replyTo,
		vlm.content,
		vlm.attachment,
		vlm.attachmentExt,
		vlm.uUuid, vlm.username, vlm.iUuid, vlm.extension,
		CASE WHEN um.status = 3 THEN true ELSE false END status
	FROM ViewLatestMessages vlm, User_Message um
	WHERE vlm.messageId = um.message AND vlm.id = ? AND um.user = ?;
`
