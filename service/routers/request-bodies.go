package routers

import "github.com/ggicci/httpin"

type UsernameRequestBody struct {
	Name string `json:"name" validate:"required,min=3,max=16"`
}

type UserPhotoRequestBody struct {
	Photo *httpin.File `in:"form=image;required"`
}

type ConversationNameRequestBody struct {
	Name string `json:"name" validate:"required,min=3,max=16"`
}

type CommentRequestBody struct {
	Comment string `json:"comment" validate:"required,length=1"`
}

type ForwardRequestBody struct {
	ForwardToId int64 `json:"destConversationId" validate:"required,min=0"`
}

type NewChatRequestBody struct {
	Recipient string `json:"recipient" validate:"required,uuid4"`
}

type NewGroupRequestBody struct {
	Name         string       `in:"form=name;required" validate:"min=3,max=16"`
	Photo        *httpin.File `in:"form=image"`
	Participants []string     `in:"form=participants" validate:"min=0,max=100"`
}

type GroupPhotoRequestBody struct {
	Photo *httpin.File `in:"form=image;required"`
}

type AddParticipantsRequestBody struct {
	Participants []string `json:"participants" validate:"required,unique,min=1,max=100"`
}

type NewMessageRequestBody struct {
	Attachment *httpin.File `in:"form=attachment" validate:"required_without=Content"`
	Content    *string      `in:"form=content" validate:"omitnil,min=0,max=4096,required_without=Attachment"`
	ReplyTo    *int64       `in:"form=repliedMessageId" validate:"omitnil,min=0"`
}
