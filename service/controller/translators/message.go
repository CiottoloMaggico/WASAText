package translators

import (
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/views"
)

func MessageInfoToCommentView(info database.MessageInfo) views.CommentView {
	return views.CommentView{
		info.GetMessageID(),
		info.GetUserUUID(),
		info.GetComment(),
	}
}

func MessageInfoListToCommentView(infos database.MessageInfoList) []views.CommentView {
	res := make([]views.CommentView, len(infos))
	for _, comment := range infos {
		res = append(res, MessageInfoToCommentView(comment))
	}
	return res
}

func MessageWithAuthorAndAttachmentToView(message database.MessageWithAuthorAndAttachment) views.MessageView {
	view := views.MessageView{
		message.GetId(),
		message.GetConversationId(),
		UserWithImageToView(message.GetAuthor()),
		message.GetStatus(),
		message.GetReplyTo(),
		nil,
		message.GetContent(),
	}

	if attach := message.GetAttachment(); attach != nil {
		imageView := ImageToView(*attach)
		view.Attachment = &imageView
	}

	return view
}

func MessageWithAuthorAndAttachmentListToView(messages database.MessageWithAuthorAndAttachmentList) []views.MessageView {
	var res = make([]views.MessageView, len(messages))
	for _, message := range messages {
		res = append(res, MessageWithAuthorAndAttachmentToView(message))
	}
	return res
}
