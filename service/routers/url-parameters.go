package routers

const DefaultPageSize = 20

type UserUrlParams struct {
	UserUUID string `url:"userUUID" validate:"required,uuid"`
}

type UserConversationUrlParams struct {
	UserUUID       string `url:"userUUID" validate:"required,uuid"`
	ConversationId int64  `url:"conversationId" validate:"required,min=0"`
}

type UserConversationMessageUrlParams struct {
	UserUUID       string `url:"userUUID" validate:"required,uuid"`
	ConversationId int64  `url:"conversationId" validate:"required,min=0"`
	MessageId      int64  `url:"messageId" validate:"required,min=0"`
}
