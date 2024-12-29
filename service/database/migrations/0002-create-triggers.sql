CREATE TRIGGER IF NOT EXISTS bind_participants_to_message
	AFTER INSERT
	ON Message
BEGIN
	INSERT INTO User_Message (message, user, status)
	SELECT new.id,
		   uc.user,
		   CASE
			   WHEN uc.user = new.author THEN 3
			   ELSE 1
			   END
	FROM User_Conversation uc
	WHERE new.conversation = uc.conversation;
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
	SELECT RAISE(ABORT, 'groups can handle at most 200 participants');
END;

