package database

const qCreateMessage = `
	INSERT INTO Message (conversation, author, replyTo, content, attachment)
	VALUES (?, ?, ?, ?, ?)
	RETURNING id, sendAt;
`

const qDeleteMessage = `
	DELETE FROM Message WHERE id = ?;
`
const qGetConversationMessageById = `
	SELECT
		id, sendAt, deliveredAt, seenAt, replyTo, content, attachment, attachmentExt, uUuid, username, iUuid, extension
	FROM ViewMessages
	WHERE id = ? AND conversation = ?;
`

const qGetConversationMessagesPaginated = `
	SELECT
		id, sendAt, deliveredAt, seenAt, replyTo, content, attachment, attachmentExt, uUuid, username, iUuid, extension
	FROM ViewMessages
	WHERE conversation = ?
	ORDER BY sendAt DESC
	LIMIT ? OFFSET ?;
`

const qGetUnseenMessages = `
	SELECT m.*
	FROM Message m, User_Message AS um
	WHERE
		um.user = ? AND um.message = id AND
		m.conversation = ? AND um.status != 3
`

const qSetSeen = `
	WITH UnSeenMessages AS ( ` + qGetUnseenMessages + ` )
	UPDATE User_Message
	SET status = 3
	FROM UnSeenMessages usm
	WHERE usm.id = id
`

const qCommentMessage = `
	UPDATE User_Message SET comment = ? WHERE message = ? AND user = ?;
`

const qDeleteComment = `
	UPDATE User_Message SET comment = NULL WHERE message = ? AND user = ?;
`
