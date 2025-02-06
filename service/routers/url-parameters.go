package routers

const DefaultPageSize = 20

type UserUrlParams struct {
	UserUUID string `url:"user_uuid" validate:"required,uuid"`
}

type UserConversationUrlParams struct {
	UserUUID       string `url:"user_uuid" validate:"required,uuid"`
	ConversationId int64  `url:"conv_id" validate:"required,min=0"`
}

type UserConversationMessageUrlParams struct {
	UserUUID       string `url:"user_uuid" validate:"required,uuid"`
	ConversationId int64  `url:"conv_id" validate:"required,min=0"`
	MessageId      int64  `url:"mess_id" validate:"required,min=0"`
}
