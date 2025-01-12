package filters

import "github.com/ciottolomaggico/wasatext/service/api/filter"

type ConversationFilterMap struct {
	Name                string `filter:"in=name,out=userConversation_name"`
	Type                string `filter:"in=type,out=userConversation_type"`
	Read                bool   `filter:"in=read,out=message_status"`
	LatestMessageSendAt string `filter:"in=latestMessage__sendAt,out=message_sendAt"`
}

func NewConversationFilter() (filter.SqlFilter, error) {
	return filter.NewSqlFilter(ConversationFilterMap{})
}
