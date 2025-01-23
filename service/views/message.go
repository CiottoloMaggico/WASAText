package views

type MessageSummaryView struct {
	Id         int64    `json:"id"`
	Author     UserView `json:"author"`
	SendAt     string   `json:"sendAt"`
	Status     string   `json:"status"`
	Content    *string  `json:"content"`
	Attachment *string  `json:"attachment"`
}

type MessageView struct {
	Id           int64             `json:"id"`
	Conversation int64             `json:"conversationId"`
	Author       UserWithImageView `json:"author"`
	SendAt       string            `json:"sendAt"`
	Status       string            `json:"status"`
	ReplyTo      *int64            `json:"repliedMessageId"`
	Attachment   *ImageView        `json:"attachment"`
	Content      *string           `json:"content"`
}

type MessageWithAuthorView struct {
	Id           int64    `json:"id"`
	Conversation int64    `json:"conversationId"`
	Author       UserView `json:"author"`
	SendAt       string   `json:"sendAt"`
	Status       string   `json:"status"`
	ReplyTo      *int64   `json:"repliedMessageId"`
	Attachment   *string  `json:"attachment"`
	Content      *string  `json:"content"`
}

type CommentView struct {
	MessageId  int64   `json:"messageId"`
	AuthorUUID string  `json:"authorUuid"`
	Content    *string `json:"content"`
}

type CommentWithAuthorView struct {
	MessageId int64             `json:"messageId"`
	Author    UserWithImageView `json:"author"`
	Content   *string           `json:"content"`
}
