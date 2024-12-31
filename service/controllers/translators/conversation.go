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
		UserConversationMessagePreviewToView(userConversation.UserConversationMessagePreview),
	}
}

func UserConversationListToView(userConversations []models.UserConversation) []views.UserConversationView {
	var res = make([]views.UserConversationView, 0, cap(userConversations))
	for _, conv := range userConversations {
		res = append(res, UserConversationToView(conv))
	}
	return res
}
