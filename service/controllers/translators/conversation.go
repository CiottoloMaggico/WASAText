package translators

import (
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
)

func UserConversationToView(userConversation models.UserConversation) views.UserConversationView {
	return views.UserConversationView{
		userConversation.Id,
		userConversation.Name,
		ImageToView(userConversation.Image),
		userConversation.Type,
		userConversation.Read,
		MessageWithAuthorToView(userConversation.MessageWithAuthor),
	}
}

func UserConversationListToView(userConversations []models.UserConversation) []views.UserConversationView {
	var res = make([]views.UserConversationView, len(userConversations))
	for _, message := range userConversations {
		res = append(res, UserConversationToView(message))
	}
	return res
}
