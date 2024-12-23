package views

type ImageView struct {
	Uuid    string `json:"uuid"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	FullUrl string `json:"fullUrl"`
}

type UserView struct {
	Uuid     string    `json:"uuid"`
	Username string    `json:"username"`
	Photo    ImageView `json:"photo"`
}

type MessageView struct {
	Id           int64      `json:"id"`
	Conversation int64      `json:"conversationId"`
	AuthorUUID   UserView   `json:"authorUUID"`
	Status       string     `json:"status"`
	ReplyTo      *int64     `json:"repliedMessageId"`
	Attachment   *ImageView `json:"attachment"`
	Content      *string    `json:"content"`
}

type CommentView struct {
	MessageId  int64   `json:"messageId"`
	AuthorUUID string  `json:"authorUUID"`
	Content    *string `json:"content"`
}
