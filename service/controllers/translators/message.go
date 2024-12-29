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
	res := make([]views.CommentView, len(infos))
	for _, comment := range infos {
		res = append(res, MessageInfoToCommentView(comment))
	}
	return res
}

func MessageWithAuthorAndAttachmentToView(message models.MessageWithAuthorAndAttachment) views.MessageView {
	view := views.MessageView{
		message.Id,
		message.Conversation,
		UserWithImageToView(message.UserWithImage),
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
	var res = make([]views.MessageView, len(messages))
	for _, message := range messages {
		res = append(res, MessageWithAuthorAndAttachmentToView(message))
	}
	return res
}

func MessageWithAuthorToView(message models.MessageWithAuthor) views.MessageWithAuthorView {
	return views.MessageWithAuthorView{
		message.Id,
		message.Conversation,
		views.PlainUser{message.User.Uuid, message.User.Uuid, message.User.Photo},
		message.GetStatus(),
		message.ReplyTo,
		message.Attachment,
		message.Content,
	}
}
