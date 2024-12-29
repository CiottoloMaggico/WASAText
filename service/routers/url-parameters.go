package routers

const DefaultPageSize = 20

type PaginationQueryParams struct {
	Page uint `validate:"required,min=0"`
	Size uint `validate:"required,min=1,max=20"`
}

type UserUrlParams struct {
	UserUUID string `validate:"required,uuid"`
}

type UserConversationUrlParams struct {
	UserUUID       string `validate:"required,uuid"`
	ConversationId int64  `validate:"required,min=0"`
}

type UserConversationMessageUrlParams struct {
	UserUUID       string `validate:"required,uuid"`
	ConversationId int64  `validate:"required,min=0"`
	MessageId      int64  `validate:"required,min=0"`
}
