package database

var tables = map[string]string{
	"MessageStatus":     qCreateTableMessageStatus,
	"Image":             qCreateTableImage,
	"User":              qCreateTableUser,
	"Conversation":      qCreateTableConversation,
	"User_Conversation": qCreateTableUser_Conversation,
	"GroupConversation": qCreateTableGroupConversation,
	"Chat":              qCreateTableChat,
	"Message":           qCreateTableMessage,
	"User_Message":      qCreateTableUser_Message,
}

const qCreateTableMessageStatus = `
	CREATE TABLE IF NOT EXISTS MessageStatus (
		id INTEGER PRIMARY KEY,
		name varchar UNIQUE NOT NULL
	);
`

const qCreateTableImage = `
	CREATE TABLE IF NOT EXISTS Image
	(
		filename varchar PRIMARY KEY,
		size     integer NOT NULL,
		owner    varchar,
		width    integer,
		height   integer,
		FOREIGN KEY (owner) REFERENCES User(uuid) ON DELETE CASCADE,
		CHECK (length(filename) >= 4 AND length(filename) <= 50),
		CHECK (size > 0 AND size <= 4096000000)
	);
`

const qCreateTableUser = `
	CREATE TABLE IF NOT EXISTS User
	(
		uuid varchar PRIMARY KEY NOT NULL, -- random uuid
		username      varchar NOT NULL,
		photoFilename varchar NOT NULL DEFAULT 'default_user_image.jpg',
		UNIQUE (username),
		FOREIGN KEY (photoFilename) REFERENCES Image (filename)
			ON DELETE RESTRICT,
		CHECK (length(username) >= 3 AND length(username) <= 16)
	);
`

const qCreateTableConversation = `
	CREATE TABLE IF NOT EXISTS Conversation
	(
		id INTEGER PRIMARY KEY,
		FOREIGN KEY (id) REFERENCES User_Conversation(conversation)
	);
`

const qCreateTableUser_Conversation = `
	CREATE TABLE IF NOT EXISTS User_Conversation
	(
		user varchar NOT NULL,
		conversation varchar NOT NULL,
		PRIMARY KEY (user, conversation),
		FOREIGN KEY (user) REFERENCES User(uuid) ON DELETE CASCADE,
		FOREIGN KEY (conversation) REFERENCES Conversation(id) ON DELETE CASCADE
	);
`

const qCreateTableChat = `
	CREATE TABLE IF NOT EXISTS Chat
	(
		id integer NOT NULL PRIMARY KEY,
		user1 varchar NOT NULL,
		user2 varchar NOT NULL,
		UNIQUE (user1, user2),
		FOREIGN KEY (id) REFERENCES Conversation(id) ON DELETE CASCADE,
		FOREIGN KEY (user1) REFERENCES User(uuid),
		FOREIGN KEY (user2) REFERENCES User(uuid)
	);
`

const qCreateTableGroupConversation = `
	CREATE TABLE IF NOT EXISTS GroupConversation (
		id integer PRIMARY KEY NOT NULL,
		name varchar NOT NULL,
		photoFilename varchar NOT NULL DEFAULT 'default_group_image.jpg',
		FOREIGN KEY (id) REFERENCES Conversation(id) ON DELETE CASCADE,
		FOREIGN KEY (photoFilename) REFERENCES Image(filename) ON DELETE RESTRICT,
		CHECK (length(name) >= 3 AND length(name) <= 16)
	);
`

const qCreateTableMessage = `
	CREATE TABLE IF NOT EXISTS Message (
		id INTEGER PRIMARY KEY,
		conversation integer NOT NULL,
		author varchar NOT NULL,
		sendAt datetime NOT NULL,
		replyTo integer,
		content varchar,
		attachmentFilename varchar,
		FOREIGN KEY (id) REFERENCES User_Message(message),
		FOREIGN KEY (conversation) REFERENCES Conversation(id) ON DELETE CASCADE,
		FOREIGN KEY (author) REFERENCES User(uuid) ON DELETE CASCADE,
		FOREIGN KEY (replyTo) REFERENCES Message(id) ON DELETE SET NULL,
		FOREIGN KEY (attachmentFilename) REFERENCES Image(filename) ON DELETE RESTRICT,
		CHECK (content IS NOT NULL OR attachmentFilename IS NOT NULL),
	    CHECK ((length(content) >= 1 AND length(content) <= 4096) OR content IS NULL)
	);
`

const qCreateTableUser_Message = `
	CREATE TABLE IF NOT EXISTS User_Message (
		message integer NOT NULL,
		user varchar NOT NULL,
		status integer NOT NULL,
		comment varchar,
		PRIMARY KEY (message, user),
		FOREIGN KEY (message) REFERENCES Message(id) ON DELETE CASCADE,
		FOREIGN KEY (user) REFERENCES User(uuid) ON DELETE CASCADE,
		FOREIGN KEY (status) REFERENCES MessageStatus (id),
		CHECK (length(comment) == 1 OR comment IS NULL)
	);
`

const qCreateViewLatestMessages = `
	CREATE VIEW IF NOT EXISTS ViewLatestMessages AS
		SELECT
			id,
			conversation,
			author,
			MAX(sendAt) as sendAt,
			replyTo,
			content,
			attachmentFilename
		FROM Message,
		GROUP BY conversation;
`
