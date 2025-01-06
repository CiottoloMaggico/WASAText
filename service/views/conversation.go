package views

type UserConversationView struct {
	Id            int64               `json:"id"`
	Name          string              `json:"name"`
	Image         ImageView           `json:"image"`
	Type          string              `json:"type"`
	Read          bool                `json:"read"`
	LatestMessage *MessageSummaryView `json:"latestMessage"`
	Participants  []UserWithImageView `json:"participants"`
}

type UserConversationSummaryView struct {
	Id            int64               `json:"id"`
	Name          string              `json:"name"`
	Image         ImageView           `json:"image"`
	Type          string              `json:"type"`
	Read          bool                `json:"read"`
	LatestMessage *MessageSummaryView `json:"latestMessage"`
}
