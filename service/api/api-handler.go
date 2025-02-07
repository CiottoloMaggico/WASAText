package api

import "net/http"

func (r *_router) Handler() http.Handler {
	r.router.POST("/session", r.wrap(r.app.GetEndpointHandler("doLogin")))
	r.router.GET("/users", r.wrap(r.app.GetEndpointHandler("getUsers")))
	r.router.GET("/users/:user_uuid", r.wrap(r.app.GetEndpointHandler("getUser")))
	r.router.PUT("/users/:user_uuid/username", r.wrap(r.app.GetEndpointHandler("setMyUserName")))
	r.router.PUT("/users/:user_uuid/avatar", r.wrap(r.app.GetEndpointHandler("setMyPhoto")))
	r.router.POST("/users/:user_uuid/chats", r.wrap(r.app.GetEndpointHandler("createChat")))
	r.router.POST("/users/:user_uuid/groups", r.wrap(r.app.GetEndpointHandler("createGroup")))
	r.router.PUT("/users/:user_uuid/groups/:group_id", r.wrap(r.app.GetEndpointHandler("addToGroup")))
	r.router.DELETE("/users/:user_uuid/groups/:group_id", r.wrap(r.app.GetEndpointHandler("leaveGroup")))
	r.router.PUT("/users/:user_uuid/groups/:group_id/name", r.wrap(r.app.GetEndpointHandler("setGroupName")))
	r.router.PUT("/users/:user_uuid/groups/:group_id/photo", r.wrap(r.app.GetEndpointHandler("setGroupPhoto")))
	r.router.GET("/users/:user_uuid/conversations", r.wrap(r.app.GetEndpointHandler("getMyConversations")))
	r.router.PUT("/users/:user_uuid/conversations", r.wrap(r.app.GetEndpointHandler("setDelivered")))
	r.router.GET("/users/:user_uuid/conversations/:conv_id", r.wrap(r.app.GetEndpointHandler("getConversation")))
	r.router.POST("/users/:user_uuid/conversations/:conv_id/messages", r.wrap(r.app.GetEndpointHandler("sendMessage")))
	r.router.PUT("/users/:user_uuid/conversations/:conv_id/messages", r.wrap(r.app.GetEndpointHandler("setSeen")))
	r.router.GET("/users/:user_uuid/conversations/:conv_id/messages", r.wrap(r.app.GetEndpointHandler("getConversationMessages")))
	r.router.GET("/users/:user_uuid/conversations/:conv_id/messages/:mess_id", r.wrap(r.app.GetEndpointHandler("messageDetail")))
	r.router.DELETE("/users/:user_uuid/conversations/:conv_id/messages/:mess_id", r.wrap(r.app.GetEndpointHandler("deleteMessage")))
	r.router.GET("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", r.wrap(r.app.GetEndpointHandler("getMessageComments")))
	r.router.PUT("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", r.wrap(r.app.GetEndpointHandler("commentMessage")))
	r.router.DELETE("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", r.wrap(r.app.GetEndpointHandler("uncommentMessage")))
	r.router.POST("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/forward", r.wrap(r.app.GetEndpointHandler("forwardMessage")))

	return r.router
}
