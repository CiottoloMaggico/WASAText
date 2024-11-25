package database

var initializers = map[string]string{
	"MessageStatus":          qSetupMessageStatus,
	"Image":                  qSetupDefaultImages,
	"User":                   "",
	"GroupConversation":      "",
	"GroupParticipants":      "",
	"GroupMessage":           "",
	"GroupMessageRecipients": "",
	"Chat":                   "",
	"ChatMessage":            "",
}

const qSetupMessageStatus = `
	INSERT INTO MessageStatus (id, name) VALUES (1, 'sent'), (2, 'delivered'), (3, 'seen');
`

const qSetupDefaultImages = `
	INSERT INTO Image (filename, size, owner, width, height)
	VALUES ('default_group_image.jpg', 135950, NULL, 840, 880),
	       ('default_user_image.jpg', 1470000, NULL, 8000, 8000);
`
