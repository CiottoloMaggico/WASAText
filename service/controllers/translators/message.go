package translators

import (
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
)

func MessageInfoToCommentView(info models.MessageInfo) views.CommentView {
	return views.CommentView{
		info.Message,
		info.User,
		info.Comment,
	}
}

func MessageInfoListToCommentView(infos []models.MessageInfo) []views.CommentView {
	res := make([]views.CommentView, 0, cap(infos))
	for _, comment := range infos {
		res = append(res, MessageInfoToCommentView(comment))
	}
	return res
}

func MessageInfoWithUserToCommentView(info models.MessageInfoWithUser) views.CommentWithAuthorView {
	return views.CommentWithAuthorView{
		info.Message,
		UserWithImageToView(info.UserWithImage),
		info.Comment,
	}
}

func MessageInfoWithUserListToCommentView(infos []models.MessageInfoWithUser) []views.CommentWithAuthorView {
	res := make([]views.CommentWithAuthorView, 0, cap(infos))
	for _, comment := range infos {
		res = append(res, MessageInfoWithUserToCommentView(comment))
	}
	return res
}

func MessageWithAuthorAndAttachmentToView(message models.MessageWithAuthorAndAttachment) views.MessageView {
	view := views.MessageView{
		message.Id,
		message.Conversation,
		UserWithImageToView(message.UserWithImage),
		message.SendAt,
		message.GetStatus(),
		message.ReplyTo,
		nil,
		message.Content,
	}

	if message.AttachmentUuid != nil {
		imageView := views.ImageView{
			*message.AttachmentUuid,
			*message.AttachmentWidth,
			*message.AttachmentHeight,
			*message.AttachmentFullUrl,
		}
		view.Attachment = &imageView
	}

	return view
}

func MessageWithAuthorAndAttachmentListToView(messages []models.MessageWithAuthorAndAttachment) []views.MessageView {
	var res = make([]views.MessageView, 0, cap(messages))
	for _, message := range messages {
		res = append(res, MessageWithAuthorAndAttachmentToView(message))
	}
	return res
}

func MessageWithAuthorToView(message models.MessageWithAuthor) views.MessageWithAuthorView {
	return views.MessageWithAuthorView{
		message.Id,
		message.Conversation,
		views.UserView{
			message.Uuid,
			message.Username,
			message.Photo,
		},
		message.SendAt,
		message.GetStatus(),
		message.ReplyTo,
		message.Attachment,
		message.Content,
	}
}

func UserConversationMessagePreviewToView(message models.UserConversationMessagePreview) *views.MessageSummaryView {
	if message.Id == nil {
		return nil
	}

	return &views.MessageSummaryView{
		*message.Id,
		views.UserView{
			*message.Uuid,
			*message.Username,
			*message.Photo,
		},
		*message.SendAt,
		message.Content,
		message.Attachment,
	}
}
