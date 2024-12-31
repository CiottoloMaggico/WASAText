package routers

import (
	"mime/multipart"
)

type UsernameRequestBody struct {
	Name string `json:"name" validate:"required,min=3,max=16"`
}

type UserPhotoRequestBody struct {
	Photo *multipart.FileHeader `in:"photo"`
}

type ConversationNameRequestBody struct {
	Name string `json:"name" validate:"required,min=3,max=16"`
}

type CommentRequestBody struct {
	Comment string `json:"comment" validate:"required"`
}

type ForwardRequestBody struct {
	ForwardToId int64 `json:"destConversationId" validate:"required,min=0"`
}

type NewChatRequestBody struct {
	Recipient string `json:"recipient" validate:"required,uuid4"`
}

type NewGroupRequestBody struct {
	Name         string                `in:"name" validate:"required,min=3,max=16"`
	Photo        *multipart.FileHeader `in:"image"`
	Participants []string              `in:"participants" validate:"min=0,max=100"`
}

type GroupPhotoRequestBody struct {
	Photo *multipart.FileHeader `in:"image"`
}

type AddParticipantsRequestBody struct {
	Participants []string `json:"participants" validate:"required,unique,min=1,max=100"`
}

type NewMessageRequestBody struct {
	Attachment *multipart.FileHeader `in:"attachment" validate:"required_without=Content"`
	Content    *string               `in:"content" validate:"omitnil,min=0,max=4096,required_without=Attachment"`
	ReplyTo    *int64                `in:"repliedMessageId" validate:"omitnil,min=0"`
}
