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
	"ViewUsers":              "",
	"ViewLatestMessages":     "",
}

const qSetupMessageStatus = `
	INSERT INTO MessageStatus (id, name) VALUES (1, 'sent'), (2, 'delivered'), (3, 'seen');
`

const qSetupDefaultImages = `
	INSERT INTO Image (uuid, extension)
	VALUES ('default_group_image', '.jpg'),
	       ('default_user_image', '.jpg');
`
