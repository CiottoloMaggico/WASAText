package database

var tables = map[string]string{
	"ForeignKey":         qEnableForeignKeyConstraints,
	"MessageStatus":      qCreateTableMessageStatus,
	"Image":              qCreateTableImage,
	"User":               qCreateTableUser,
	"Conversation":       qCreateTableConversation,
	"User_Conversation":  qCreateTableUser_Conversation,
	"GroupConversation":  qCreateTableGroupConversation,
	"Chat":               qCreateTableChat,
	"Message":            qCreateTableMessage,
	"User_Message":       qCreateTableUser_Message,
	"ViewUsers":          qCreateViewUsers,
	"ViewMessages":       qCreateViewMessages,
	"ViewLatestMessages": qCreateViewLatestMessages,
	"ViewConversations":  qCreateViewConversation,
}

const qEnableForeignKeyConstraints = `PRAGMA foreign_keys = ON;`

const qCreateTableMessageStatus = `
	CREATE TABLE IF NOT EXISTS MessageStatus (
		id INTEGER PRIMARY KEY,
		name varchar UNIQUE NOT NULL
	);
`

const qCreateTableImage = `
	CREATE TABLE IF NOT EXISTS Image
	(
		uuid  PRIMARY KEY,
		extension varchar NOT NULL,
		uploadedAt datetime NOT NULL DEFAULT current_timestamp
	);
`

const qCreateTableUser = `
	CREATE TABLE IF NOT EXISTS User
	(
		uuid varchar PRIMARY KEY NOT NULL, -- random uuid
		username      varchar NOT NULL,
		photo varchar NOT NULL DEFAULT 'default_user_image',
		UNIQUE (username),
		FOREIGN KEY (photo) REFERENCES Image (uuid)
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
		author varchar NOT NULL,
		photo varchar NOT NULL DEFAULT 'default_group_image',
		FOREIGN KEY (id) REFERENCES Conversation(id) ON DELETE CASCADE,
		FOREIGN KEY (author) REFERENCES User(uuid),
		FOREIGN KEY (photo) REFERENCES Image(uuid) ON DELETE RESTRICT,
		CHECK (length(name) >= 3 AND length(name) <= 16)
	);
`

const qCreateTableMessage = `
	CREATE TABLE IF NOT EXISTS Message (
		id INTEGER PRIMARY KEY,
		conversation integer NOT NULL,
		author varchar NOT NULL,
		sendAt datetime NOT NULL DEFAULT current_timestamp,
		deliveredAt datetime,
		seenAt datetime,
		replyTo integer,
		content varchar,
		attachment varchar,
		FOREIGN KEY (id) REFERENCES User_Message(message),
		FOREIGN KEY (conversation) REFERENCES Conversation(id) ON DELETE CASCADE,
		FOREIGN KEY (author) REFERENCES User(uuid) ON DELETE CASCADE,
		FOREIGN KEY (replyTo) REFERENCES Message(id) ON DELETE SET NULL,
		FOREIGN KEY (attachment) REFERENCES Image(uuid) ON DELETE RESTRICT,
		CHECK (content IS NOT NULL OR attachment IS NOT NULL),
	    CHECK ((length(content) >= 1 AND length(content) <= 4096) OR content IS NULL)
	);
`

const qCreateTableUser_Message = `
	CREATE TABLE IF NOT EXISTS User_Message (
		message integer NOT NULL,
		user varchar NOT NULL,
		status integer NOT NULL DEFAULT 1,
		comment varchar,
		PRIMARY KEY (message, user),
		FOREIGN KEY (message) REFERENCES Message(id) ON DELETE CASCADE,
		FOREIGN KEY (user) REFERENCES User(uuid) ON DELETE CASCADE,
		FOREIGN KEY (status) REFERENCES MessageStatus (id),
		CHECK (length(comment) == 1 OR comment IS NULL)
	);
`

const qCreateViewUsers = `
	CREATE VIEW IF NOT EXISTS ViewUsers AS
		SELECT User.uuid AS uUuid, username, Image.uuid AS iUuid, extension
		FROM User, Image
		WHERE User.photo = Image.uuid;
`

const qCreateViewMessages = `
	CREATE VIEW IF NOT EXISTS ViewMessages AS
		SELECT
			m.id, m.conversation, m.sendAt, m.deliveredAt, m.seenAt, m.replyTo, m.content, m.attachment, i.extension attachmentExt,
			u.*
		FROM Message m, ViewUsers u
		LEFT OUTER JOIN Image i ON i.uuid = m.attachment
		WHERE
			m.author = u.uUuid;
`

const qCreateViewLatestMessages = `
	CREATE VIEW IF NOT EXISTS ViewLatestMessages AS
		SELECT
			c.id,
			vm.id messageId,
			MAX(vm.sendAt) sendAt,
			deliveredAt,
			seenAt,
			vm.replyTo,
			vm.content,
			vm.attachment,
			vm.attachmentExt,
			vm.uUuid, vm.username, vm.iUuid, vm.extension
		FROM Conversation c
		LEFT OUTER JOIN ViewMessages vm ON vm.conversation = c.id
		GROUP BY c.id;
`

const qCreateViewConversation = `
	CREATE VIEW IF NOT EXISTS ViewConversations AS
	SELECT
		Conversation.id id,
		vc.user1,
		vc.user2,
		vg.name,
		vg.photo
	FROM Conversation
	LEFT OUTER JOIN Chat vc ON (vc.id = Conversation.id)
    LEFT OUTER JOIN GroupConversation vg ON (vg.id = Conversation.id)
`
