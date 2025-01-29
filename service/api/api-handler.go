package api

import "net/http"

func (r *_router) Handler() http.Handler {
	r.router.POST("/session", r.wrap(r.doLogin))
	r.router.GET("/users", r.wrap(r.authMiddleware.Wrap(r.getUsers)))
	r.router.GET("/users/:userUUID", r.wrap(r.authMiddleware.Wrap(r.getUser)))
	r.router.PUT("/users/:userUUID/username", r.wrap(r.authMiddleware.Wrap(r.setMyUserName)))
	r.router.PUT("/users/:userUUID/avatar", r.wrap(r.authMiddleware.Wrap(r.setMyPhoto)))
	r.router.POST("/users/:userUUID/chats", r.wrap(r.authMiddleware.Wrap(r.createChat)))
	r.router.POST("/users/:userUUID/groups", r.wrap(r.authMiddleware.Wrap(r.createGroup)))
	r.router.PUT("/users/:userUUID/groups/:groupId", r.wrap(r.authMiddleware.Wrap(r.addToGroup)))
	r.router.DELETE("/users/:userUUID/groups/:groupId", r.wrap(r.authMiddleware.Wrap(r.leaveGroup)))
	r.router.PUT("/users/:userUUID/groups/:groupId/name", r.wrap(r.authMiddleware.Wrap(r.setGroupName)))
	r.router.PUT("/users/:userUUID/groups/:groupId/photo", r.wrap(r.authMiddleware.Wrap(r.setGroupPhoto)))
	r.router.GET("/users/:userUUID/conversations", r.wrap(r.authMiddleware.Wrap(r.getMyConversations)))
	r.router.PUT("/users/:userUUID/conversations", r.wrap(r.authMiddleware.Wrap(r.setDelivered)))
	r.router.GET("/users/:userUUID/conversations/:convId", r.wrap(r.authMiddleware.Wrap(r.getConversation)))
	r.router.POST("/users/:userUUID/conversations/:convId/messages", r.wrap(r.authMiddleware.Wrap(r.sendMessage)))
	r.router.PUT("/users/:userUUID/conversations/:convId/messages", r.wrap(r.authMiddleware.Wrap(r.setSeen)))
	r.router.GET("/users/:userUUID/conversations/:convId/messages", r.wrap(r.authMiddleware.Wrap(r.getConversationMessages)))
	r.router.GET("/users/:userUUID/conversations/:convId/messages/:messId", r.wrap(r.authMiddleware.Wrap(r.messageDetail)))
	r.router.DELETE("/users/:userUUID/conversations/:convId/messages/:messId", r.wrap(r.authMiddleware.Wrap(r.deleteMessage)))
	r.router.GET("/users/:userUUID/conversations/:convId/messages/:messId/comments", r.wrap(r.authMiddleware.Wrap(r.getMessageComments)))
	r.router.PUT("/users/:userUUID/conversations/:convId/messages/:messId/comments", r.wrap(r.authMiddleware.Wrap(r.commentMessage)))
	r.router.DELETE("/users/:userUUID/conversations/:convId/messages/:messId/comments", r.wrap(r.authMiddleware.Wrap(r.uncommentMessage)))
	r.router.POST("/users/:userUUID/conversations/:convId/messages/:messId/forward", r.wrap(r.authMiddleware.Wrap(r.forwardMessage)))

	return r.router
}
