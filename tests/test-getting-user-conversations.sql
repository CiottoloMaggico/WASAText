WITH LatestMessages AS (SELECT id,
							   conversation,
							   author,
							   MAX(sendAt) as sendAt,
							   status,
							   replyTo,
							   content,
							   attachmentFilename
						FROM Message,
							 User_Message
						WHERE User_Message.message = Message.id
						  AND User_Message.user = ?
						GROUP BY conversation),
	 Conversations AS (SELECT chat.id            AS id,
							  pers.username      AS name,
							  pers.photoFilename AS photoFilename
					   FROM Chat chat,
							User pers
					   WHERE (chat.user1 = ? AND pers.uuid = chat.user2)
						  OR (chat.user2 = ? AND pers.uuid = chat.user1)
					   UNION
					   SELECT g.id            AS id,
							  g.name          AS name,
							  g.photoFilename AS photoFilename
					   FROM User_Conversation uc,
							GroupConversation g
					   WHERE uc.conversation = g.id
						 AND uc.user = ?)
SELECT conv.id,
	   conv.name,
	   conv.photoFilename,
	   mess.*
FROM Conversations AS conv
		 LEFT OUTER JOIN
	 LatestMessages mess
	 ON mess.conversation = conv.id
ORDER BY mess.sendAt DESC
