package translators

import (
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
)

func UserConversationToSummaryView(userConversation models.UserConversation) views.UserConversationSummaryView {
	return views.UserConversationSummaryView{
		userConversation.Id,
		userConversation.Name,
		ImageToView(userConversation.Image),
		userConversation.Type,
		userConversation.Read,
		UserConversationMessagePreviewToView(userConversation.UserConversationMessagePreview),
	}
}

func UserConversationListToSummaryView(userConversations []models.UserConversation) []views.UserConversationSummaryView {
	var res = make([]views.UserConversationSummaryView, 0, cap(userConversations))
	for _, conv := range userConversations {
		res = append(res, UserConversationToSummaryView(conv))
	}
	return res
}

func UserConversationToView(userConversation models.UserConversation, participants []models.UserWithImage) views.UserConversationView {
	return views.UserConversationView{
		userConversation.Id,
		userConversation.Name,
		ImageToView(userConversation.Image),
		userConversation.Type,
		userConversation.Read,
		UserConversationMessagePreviewToView(userConversation.UserConversationMessagePreview),
		UserWithImageListToView(participants),
	}
}
