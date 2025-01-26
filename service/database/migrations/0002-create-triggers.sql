CREATE TRIGGER IF NOT EXISTS bind_participants_to_message
	AFTER INSERT
	ON Message
BEGIN
	INSERT INTO User_Message (message, user, status)
	SELECT new.id,
		   uc.user,
		   CASE WHEN uc.user = new.author THEN 3 ELSE 1 END
	FROM User_Conversation uc
	WHERE new.conversation = uc.conversation;
END;

CREATE TRIGGER IF NOT EXISTS message_author_must_be_a_participant
	BEFORE INSERT
	ON Message
	WHEN
		NOT EXISTS(SELECT *
				   FROM User_Conversation
				   WHERE conversation = new.conversation
					 AND user = new.author)
BEGIN
	SELECT RAISE(ABORT, 'TRIGGER: message_author_must_be_a_participant');
END;

CREATE TRIGGER IF NOT EXISTS replied_message_within_same_conversation
	BEFORE INSERT
	ON Message
	WHEN
		new.replyTo IS NOT NULL
			AND
		NOT EXISTS(SELECT *
				   FROM Message
				   WHERE message.id = new.replyTo
					 AND message.conversation = new.conversation)
BEGIN
	SELECT RAISE(ABORT,
				 'TRIGGER: replied_message_within_same_conversation');
END;

CREATE TRIGGER IF NOT EXISTS add_chat_to_users
	AFTER INSERT
	ON Chat
BEGIN
	INSERT INTO User_Conversation
	VALUES (new.user1, new.id),
		   (new.user2, new.id);
END;

CREATE TRIGGER IF NOT EXISTS add_group_to_author
	AFTER INSERT
	ON GroupConversation
BEGIN
	INSERT INTO User_Conversation
	VALUES (new.author, new.id);
END;

CREATE TRIGGER IF NOT EXISTS set_message_status_to_delivered
	AFTER UPDATE OF status
	ON User_Message
	WHEN
		NOT EXISTS (SELECT um.message, um.user
					FROM User_Message um
					WHERE um.message = new.message
					  AND status = 1)
			AND
		EXISTS (SELECT m.id
				FROM Message m
				WHERE m.id = new.message
				  AND deliveredAt IS NULL)
BEGIN
	UPDATE Message SET deliveredAt = current_timestamp WHERE id = new.message;
END;



CREATE TRIGGER IF NOT EXISTS set_message_status_to_seen
	AFTER UPDATE OF status
	ON User_Message
	WHEN
		NOT EXISTS (SELECT um.message, um.user
					FROM User_Message um
					WHERE um.message = new.message
					  AND (status = 1 OR status = 2))
			AND
		EXISTS (SELECT m.id
				FROM Message m
				WHERE m.id = new.message
				  AND seenAt IS NULL)
BEGIN
	UPDATE Message SET seenAt = current_timestamp WHERE id = new.message;
END;

CREATE TRIGGER IF NOT EXISTS group_participants_limit
	BEFORE INSERT
	ON User_Conversation
	WHEN
		(SELECT COUNT(*)
		 FROM User_Conversation
		 WHERE conversation = new.conversation) >= 200
BEGIN
	SELECT RAISE(ABORT, 'TRIGGER: group_participants_limit');
END;

CREATE TRIGGER IF NOT EXISTS delete_empty_groups
	AFTER DELETE
	ON User_Conversation
	WHEN
		(SELECT COUNT(*)
		 FROM User_Conversation
		 WHERE conversation = old.conversation
		   AND EXISTS(SELECT * FROM GroupConversation WHERE old.conversation = GroupConversation.id)) = 0
BEGIN
	DELETE FROM Conversation WHERE id = old.conversation;
END;

