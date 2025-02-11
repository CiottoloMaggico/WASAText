CREATE TABLE IF NOT EXISTS MessageStatus
(
	id   INTEGER PRIMARY KEY,
	name varchar UNIQUE NOT NULL
);



CREATE TABLE IF NOT EXISTS Image
(
	uuid PRIMARY KEY,
	extension  varchar  NOT NULL,
	width      integer  NOT NULL,
	height     integer  NOT NULL,
	fullUrl    varchar  NOT NULL,
	uploadedAt datetime NOT NULL DEFAULT current_timestamp
);



CREATE TABLE IF NOT EXISTS User
(
	uuid     varchar PRIMARY KEY NOT NULL, -- random uuid
	username varchar             NOT NULL,
	photo    varchar             NOT NULL DEFAULT 'default_user_image',
	UNIQUE (username),
	FOREIGN KEY (photo) REFERENCES Image (uuid)
		ON DELETE RESTRICT,
	CHECK (length(username) >= 3 AND length(username) <= 16)
);



CREATE TABLE IF NOT EXISTS Conversation
(
	id INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS User_Conversation
(
	user         varchar NOT NULL,
	conversation varchar NOT NULL,
	PRIMARY KEY (user, conversation),
	FOREIGN KEY (user) REFERENCES User (uuid) ON DELETE CASCADE,
	FOREIGN KEY (conversation) REFERENCES Conversation (id) ON DELETE CASCADE
);



CREATE TABLE IF NOT EXISTS Chat
(
	id    integer NOT NULL PRIMARY KEY,
	user1 varchar NOT NULL,
	user2 varchar NOT NULL,
	UNIQUE (user1, user2),
	FOREIGN KEY (id) REFERENCES Conversation (id) ON DELETE CASCADE,
	FOREIGN KEY (user1) REFERENCES User (uuid),
	FOREIGN KEY (user2) REFERENCES User (uuid)
);



CREATE TABLE IF NOT EXISTS GroupConversation
(
	id     integer PRIMARY KEY NOT NULL,
	name   varchar             NOT NULL,
	author varchar             NOT NULL,
	photo  varchar             NOT NULL DEFAULT 'default_group_image',
	FOREIGN KEY (id) REFERENCES Conversation (id) ON DELETE CASCADE,
	FOREIGN KEY (author) REFERENCES User (uuid),
	FOREIGN KEY (photo) REFERENCES Image (uuid) ON DELETE RESTRICT,
	CHECK (length(name) >= 3 AND length(name) <= 16)
);



CREATE TABLE IF NOT EXISTS Message
(
	id           INTEGER PRIMARY KEY,
	conversation integer  NOT NULL,
	author       varchar  NOT NULL,
	sendAt       datetime NOT NULL DEFAULT current_timestamp,
	deliveredAt  datetime,
	seenAt       datetime,
	replyTo      integer,
	content      varchar,
	attachment   varchar,
	FOREIGN KEY (conversation) REFERENCES Conversation (id) ON DELETE CASCADE,
	FOREIGN KEY (author) REFERENCES User (uuid) ON DELETE CASCADE,
	FOREIGN KEY (replyTo) REFERENCES Message (id) ON DELETE SET NULL,
	FOREIGN KEY (attachment) REFERENCES Image (uuid) ON DELETE RESTRICT,
	CHECK (content IS NOT NULL OR attachment IS NOT NULL),
	CHECK ((length(content) >= 1 AND length(content) <= 4096) OR content IS NULL)
);

CREATE TABLE IF NOT EXISTS User_Message
(
	message integer NOT NULL,
	user    varchar NOT NULL,
	status  integer NOT NULL DEFAULT 1,
	comment varchar,
	PRIMARY KEY (message, user),
	FOREIGN KEY (message) REFERENCES Message (id) ON DELETE CASCADE,
	FOREIGN KEY (user) REFERENCES User (uuid) ON DELETE CASCADE,
	FOREIGN KEY (status) REFERENCES MessageStatus (id)
);

CREATE VIEW IF NOT EXISTS ViewUsers AS
SELECT User.uuid  user_uuid,
	   username   user_username,
	   Image.uuid image_uuid,
	   extension  image_extension,
	   width      image_width,
	   height     image_height,
	   fullUrl    image_fullUrl,
	   uploadedAt image_uploadedAt
FROM User,
	 Image
WHERE User.photo = Image.uuid;

CREATE VIEW IF NOT EXISTS ViewMessages AS
SELECT m.id               message_id,
	   m.conversation     message_conversation,
	   m.sendAt           message_sendAt,
	   m.deliveredAt      message_deliveredAt,
	   m.seenAt           message_seenAt,
	   m.replyTo          message_replyTo,
	   m.content          message_content,
	   i.uuid             attachment_uuid,
	   i.extension        attachment_extension,
	   i.width            attachment_width,
	   i.height           attachment_height,
	   i.fullUrl          attachment_fullUrl,
	   i.uploadedAt       attachment_uploadedAt,
	   u.*
FROM Message m,
	 ViewUsers u
		 LEFT OUTER JOIN Image i ON i.uuid = m.attachment
WHERE m.author = u.user_uuid;

CREATE VIEW IF NOT EXISTS ViewLatestMessages AS
SELECT m.id           message_id,
	   m.conversation message_conversation,
	   MAX(m.sendAt)  message_sendAt,
	   m.deliveredAt  message_deliveredAt,
	   m.seenAt       message_seenAt,
	   m.replyTo      message_replyTo,
	   m.content      message_content,
	   i.uuid         attachment_uuid,
	   i.extension    attachment_extension,
	   i.width        attachment_width,
	   i.height       attachment_height,
	   i.fullUrl      attachment_fullUrl,
	   i.uploadedAt   attachment_uploadedAt,
	   u.*
FROM Message m,
	 ViewUsers u
		 LEFT OUTER JOIN Image i ON i.uuid = m.attachment
WHERE m.author = u.user_uuid
GROUP BY m.conversation;


