package api

import "net/http"

func (r *_router) Handler() http.Handler {
	r.router.POST("/session", r.wrap(r.doLogin))
	r.router.GET("/users", r.wrap(r.authMiddleware.Wrap(r.getUsers)))
	r.router.GET("/users/:user_uuid", r.wrap(r.authMiddleware.Wrap(r.getUser)))
	r.router.PUT("/users/:user_uuid/username", r.wrap(r.authMiddleware.Wrap(r.setMyUserName)))
	r.router.PUT("/users/:user_uuid/avatar", r.wrap(r.authMiddleware.Wrap(r.setMyPhoto)))
	r.router.POST("/users/:user_uuid/chats", r.wrap(r.authMiddleware.Wrap(r.createChat)))
	r.router.POST("/users/:user_uuid/groups", r.wrap(r.authMiddleware.Wrap(r.createGroup)))
	r.router.PUT("/users/:user_uuid/groups/:group_id", r.wrap(r.authMiddleware.Wrap(r.addToGroup)))
	r.router.DELETE("/users/:user_uuid/groups/:group_id", r.wrap(r.authMiddleware.Wrap(r.leaveGroup)))
	r.router.PUT("/users/:user_uuid/groups/:group_id/name", r.wrap(r.authMiddleware.Wrap(r.setGroupName)))
	r.router.PUT("/users/:user_uuid/groups/:group_id/photo", r.wrap(r.authMiddleware.Wrap(r.setGroupPhoto)))
	r.router.GET("/users/:user_uuid/conversations", r.wrap(r.authMiddleware.Wrap(r.getMyConversations)))
	r.router.PUT("/users/:user_uuid/conversations", r.wrap(r.authMiddleware.Wrap(r.setDelivered)))
	r.router.GET("/users/:user_uuid/conversations/:conv_id", r.wrap(r.authMiddleware.Wrap(r.getConversation)))
	r.router.POST("/users/:user_uuid/conversations/:conv_id/messages", r.wrap(r.authMiddleware.Wrap(r.sendMessage)))
	r.router.PUT("/users/:user_uuid/conversations/:conv_id/messages", r.wrap(r.authMiddleware.Wrap(r.setSeen)))
	r.router.GET("/users/:user_uuid/conversations/:conv_id/messages", r.wrap(r.authMiddleware.Wrap(r.getConversationMessages)))
	r.router.GET("/users/:user_uuid/conversations/:conv_id/messages/:mess_id", r.wrap(r.authMiddleware.Wrap(r.messageDetail)))
	r.router.DELETE("/users/:user_uuid/conversations/:conv_id/messages/:mess_id", r.wrap(r.authMiddleware.Wrap(r.deleteMessage)))
	r.router.GET("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", r.wrap(r.authMiddleware.Wrap(r.getMessageComments)))
	r.router.PUT("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", r.wrap(r.authMiddleware.Wrap(r.commentMessage)))
	r.router.DELETE("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", r.wrap(r.authMiddleware.Wrap(r.uncommentMessage)))
	r.router.POST("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/forward", r.wrap(r.authMiddleware.Wrap(r.forwardMessage)))

	return r.router
}
