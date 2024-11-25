package database

var triggers = [...]string{
	tInsertOnGroupConversationDisjoint,
	tUpdateOnGroupConversationDisjoint,
	tInsertOnChatDisjoint,
	tUpdateOnChatDisjoint,
	tInsertOnConversationComplete,
	tDeleteOnChat,
	tDeleteOnGroupConversation,
	tInsertOnChat,
	tInsertOnMessage,
}

// If the inserted or updated row breaks the disjoint property between chat and group chat raise an error
const tbInsertOrUpdateOnGroupConversationOrChat = `
	BEGIN
		SELECT
			CASE
				WHEN EXISTS(SELECT * FROM Chat AS c, GroupConversation AS g WHERE new.id = c.id AND new.id = g.id)
				THEN RAISE(ABORT, "Conversation type must be disjoint")
			END;
	END;
`

const tInsertOnGroupConversationDisjoint = `
	CREATE TRIGGER IF NOT EXISTS disjoint_type_of_conversation_group_conversation_insert AFTER INSERT ON GroupConversation FOR EACH ROW
` + tbInsertOrUpdateOnGroupConversationOrChat

const tUpdateOnGroupConversationDisjoint = `
	CREATE TRIGGER IF NOT EXISTS disjoint_type_of_conversation_group_conversation_update AFTER UPDATE ON GroupConversation FOR EACH ROW
` + tbInsertOrUpdateOnGroupConversationOrChat

const tInsertOnChatDisjoint = `
	CREATE TRIGGER IF NOT EXISTS disjoint_type_of_conversation_chat_insert AFTER INSERT ON Chat FOR EACH ROW
` + tbInsertOrUpdateOnGroupConversationOrChat

const tUpdateOnChatDisjoint = `
	CREATE TRIGGER IF NOT EXISTS disjoint_type_of_conversation_chat_update AFTER UPDATE ON Chat FOR EACH ROW
` + tbInsertOrUpdateOnGroupConversationOrChat

// If the inserted row breaks the complete property between chat and group chat raise an error
const tbInsertOnConversation = `
	BEGIN
		SELECT
			CASE
				WHEN NOT EXISTS(
						SELECT c.id FROM Chat AS c WHERE new.id = c.id
						UNION
						SELECT g.id FROM GroupConversation AS g WHERE new.id = g.id
					)
				THEN RAISE(ABORT, "Conversation type must be complete")
			END;
	END;
`

const tInsertOnConversationComplete = `
	CREATE TRIGGER IF NOT EXISTS complete_conversation_type AFTER INSERT ON Conversation FOR EACH ROW
` + tbInsertOnConversation

// After the cancellation of a chat or a group chat always delete also the conversation
const tbDeleteOnChatOrGroupConversation = `
	BEGIN
		DELETE FROM Conversation WHERE id = old.id;
	END;
`

const tDeleteOnChat = `
	CREATE TRIGGER IF NOT EXISTS delete_on_chat AFTER DELETE ON Chat FOR EACH ROW
` + tbDeleteOnChatOrGroupConversation

const tDeleteOnGroupConversation = `
	CREATE TRIGGER IF NOT EXISTS delete_on_group_conversation AFTER DELETE ON GroupConversation FOR EACH ROW
` + tbDeleteOnChatOrGroupConversation

// After the insertion of a new chat always add the chat to the conversation list of both users
const tbInsertOnChat = `
	BEGIN
		INSERT INTO User_Conversation VALUES (new.user1, new.id), (new.user2, new.id);
	END;
`

const tInsertOnChat = `
	CREATE TRIGGER IF NOT EXISTS insert_on_chat AFTER INSERT ON Chat FOR EACH ROW
` + tbInsertOnChat

// After the insertion of a new message always add the author of the message to the User_Message table
const tbInsertOnMessage = `
	BEGIN
		INSERT INTO User_Message VALUES (new.id, new.author, 'seen', NULL);
	END;
`

const tInsertOnMessage = `
	CREATE TRIGGER IF NOT EXISTS insert_on_message AFTER INSERT ON Message FOR EACH ROW
` + tbInsertOnMessage
