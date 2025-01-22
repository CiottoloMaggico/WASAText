package filters

import "github.com/ciottolomaggico/wasatext/service/api/filter"

type MessageFilterMap struct {
	AuthorUUID       string `filter:"in=author__uuid,out=user_uuid"`
	SendAt           string `filter:"in=sendAt,out=message_sendAt"`
	RepliedMessageId int    `filter:"in=repliedMessageId,out=message_replyTo"`
}

func NewMessageFilter() (filter.SqlFilter, error) {
	return filter.NewSqlFilter(MessageFilterMap{})
}
