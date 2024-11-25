package database

const qCreateMessage = `
	INSERT INTO Message (conversation, author, sendAt, replyTo, content, attachmentFilename)
	VALUES (?, ?, ?, ?, ?, ?) RETURNING id;
`
const qCreateUser_Message = `
	INSERT INTO User_Message (message, user, status)
	VALUES (?, ?, ?);
`

const qGetMessages = `
	SELECT *
	FROM Message
	WHERE conversation = ?
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

const qDeleteMessage = `
	DELETE FROM Message WHERE id = ?;
`

const qCommentMessage = `
	UPDATE User_Message SET comment = ? WHERE message = ? AND user = ?;
`

const qDeleteComment = `
	UPDATE User_Message SET comment = NULL WHERE message = ? AND user = ?;
`
