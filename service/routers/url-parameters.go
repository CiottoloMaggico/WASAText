package routers

const DefaultPageSize = 20

type UserUrlParams struct {
	UserUUID string `in:"userUUID" validate:"required,uuid"`
}

type UserConversationUrlParams struct {
	UserUUID       string `in:"userUUID" validate:"required,uuid"`
	ConversationId int64  `in:"conversationId" validate:"required,min=0"`
}

type UserConversationMessageUrlParams struct {
	UserUUID       string `in:"userUUID" validate:"required,uuid"`
	ConversationId int64  `in:"conversationId" validate:"required,min=0"`
	MessageId      int64  `in:"messageId" validate:"required,min=0"`
}
