package database

const qLatestMessagesWithStatus = `
	SELECT
		vlm.id AS id,
		vlm.conversation AS conversation,
		vlm.author AS author,
		vlm.sendAt AS sendAt,
		um.status AS status,
		vlm.replyTo AS replyTo,
		vlm.content AS content,
		vlm.attachmentFilename AS attachmentFilename
	FROM ViewLatestMessages vlm,
		User_Message um
	WHERE
		um.user = ? AND
		um.message = vlm.id
`

const qConversationLatestMessage = qLatestMessagesWithStatus + ` AND vlm.conversation = ?`

const qGetMyChats = `
	SELECT
		chat.id AS id,
		pers.username AS name,
		pers.photoFilename AS photoFilename,
		'one_to_one' AS type
	FROM Chat chat,
		User pers
	WHERE
		(chat.user1 = ? AND pers.uuid = chat.user2)
		OR
		(chat.user2 = ? AND pers.uuid = chat.user1)
`

const qGetMyGroups = `
	SELECT
		g.id AS id,
		g.name AS name,
		g.photoFilename AS photoFilename,
		'group' AS type
	FROM User_Conversation uc,
		GroupConversation g
	WHERE uc.conversation = g.id AND uc.user = ?
`

const qGetMyConversations = qGetMyChats + ` UNION ` + qGetMyGroups

const qGetConversation = qGetMyChats + ` AND chat.id = ? UNION ` + qGetMyGroups + ` AND g.id = ?`

const qGetMyConversationsWithLatestMessage = `
	WITH UserMessages AS ( ` + qLatestMessagesWithStatus + ` ),
	Conversations AS (` + qGetMyConversations + ` )
	SELECT
		conv.id,
		conv.name,
		conv.photoFilename,
		conv.type,
		mess.id,
		mess.author,
		mess.sendAt,
		mess.status,
		mess.replyTo,
		mess.content,
		mess.attachmentFilename
	FROM Conversations AS conv
	LEFT OUTER JOIN LatestMessages mess ON mess.conversation = conv.id
	ORDER BY status ASC
`

const qGetMyConversationsWithLatestMessagePaginated = qGetMyConversationsWithLatestMessage + ` LIMIT ? OFFSET ?;`

const qGetConversationWithLatestMessage = `
	WITH LatestMessage AS ( ` + qConversationLatestMessage + ` ),
	Conversation AS ( ` + qGetConversation + ` )
	SELECT
		conv.id,
		conv.name,
		conv.photoFilename,
		conv.type,
		mess.id,
		mess.author,
		mess.sendAt,
		mess.status,
		mess.replyTo,
		mess.content,
		mess.attachmentFilename
	FROM Conversations AS conv
	LEFT OUTER JOIN LatestMessages mess ON mess.conversation = conv.id;
`

const qSetDelivered = `
	UPDATE User_Message SET status = 2 WHERE user = ? AND status = 1;
`

const qGetConversationParticipants = `
	SELECT uc.user
	FROM User_Conversation uc
	WHERE uc.conversation = ?
`

const qCreateConversation = `
	INSERT INTO Conversation DEFAULT VALUES RETURNING id;
`

const qCreateChat = `
	INSERT INTO Chat(id, user1, user2) VALUES (?, ?, ?);
	INSERT INTO User_Conversation (user, conversation)
	VALUES (?, ?), (?, ?); -- (user1, id), (user2, id)

`

const qCreateGroup = `
	INSERT INTO GroupConversation (id, name, photoFilename) VALUES (?, ?, ?);
	INSERT INTO User_Conversation (user, conversation)
	VALUES (?, ?); -- (author, id)
`

const qUpdateGroupName = `
	UPDATE GroupConversation SET name = ? WHERE id = ?;
`

const qUpdateGroupImage = `
	UPDATE GroupConversation SET photoFilename = ? WHERE id = ?;
`

const qAddGroupParticipant = `
	INSERT INTO User_Conversation (user, conversation)
	VALUES (?, ?); -- (user, id)
`

const qRemoveGroupParticipant = `
	DELETE FROM User_Conversation
	WHERE user = ? AND conversation = ?
`
