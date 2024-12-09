package database

const qLinkUserToConversation = `
	INSERT INTO User_Conversation VALUES (?, ?);
`

const qUnLinkUserToConversation = `
	DELETE FROM User_Conversation WHERE user = ? AND conversation = ?;
`

const qGetConversationMessagesPaginated = `
	SELECT
		id, sendAt, deliveredAt, seenAt, replyTo, content, attachment, attachmentExt, uUuid, username, iUuid, extension
	FROM ViewMessages
	WHERE conversation = ?
	ORDER BY sendAt DESC
	LIMIT ? OFFSET ?;
`

const qGetLatestMessage = `
	SELECT
		id, MAX(sendAt) sendAt, deliveredAt, seenAt, replyTo, content, attachment, attachmentExt, uUuid, username, iUuid, extension
	FROM ViewMessages
	WHERE conversation = ?;
`

const qGetConversation = `
	SELECT
		vc.*
	FROM ViewConversations vc, User_Conversation uc
	WHERE uc.conversation = vc.id AND vc.id = ? AND uc.user = ?;
`

const qGetGroup = `
	SELECT
		id,
	    name,
	    vu.uUuid,
	    vu.username,
	    vu.iUuid,
	    vu.extension,
		photo,
		Image.extension
	FROM GroupConversation,
	     ViewUsers vu,
	     Image
	WHERE id = ? AND photo = Image.uuid AND vu.uUuid = author;
`

const qGetChat = `
	SELECT
		c.id,
		user1.uUuid uuid1,
		user1.username username1,
		user1.iUuid image1,
		user1.extension extension1,
		user2.uUuid uuid2,
		user2.username username2,
		user2.iUuid image2,
		user2.extension extension2
	FROM Chat c, ViewUsers user1, ViewUsers user2
	WHERE c.user1 = user1.uUuid AND c.user2 = user2.uUuid;

`

//OLD

const getUserConversationsPaginated = `
	SELECT
		vc.*
	FROM ViewConversations vc, User_Conversation uc
	WHERE uc.conversation = vc.id AND uc.user = ?
	LIMIT ? OFFSET ?;
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

const qGetParticipant = `
	SELECT uc.user FROM User_Conversation uc WHERE uc.user = ? AND uc.conversation = ?;
`
