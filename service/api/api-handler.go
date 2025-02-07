package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (r *_router) Handler() http.Handler {
	r.router.POST("/session", r.Handle("doLogin"))
	r.router.GET("/users", r.Handle("getUsers"))
	r.router.GET("/users/:user_uuid", r.Handle("getUser"))
	r.router.PUT("/users/:user_uuid/username", r.Handle("setMyUserName"))
	r.router.PUT("/users/:user_uuid/avatar", r.Handle("setMyPhoto"))
	r.router.POST("/users/:user_uuid/chats", r.Handle("createChat"))
	r.router.POST("/users/:user_uuid/groups", r.Handle("createGroup"))
	r.router.PUT("/users/:user_uuid/groups/:group_id", r.Handle("addToGroup"))
	r.router.DELETE("/users/:user_uuid/groups/:group_id", r.Handle("leaveGroup"))
	r.router.PUT("/users/:user_uuid/groups/:group_id/name", r.Handle("setGroupName"))
	r.router.PUT("/users/:user_uuid/groups/:group_id/photo", r.Handle("setGroupPhoto"))
	r.router.GET("/users/:user_uuid/conversations", r.Handle("getMyConversations"))
	r.router.PUT("/users/:user_uuid/conversations", r.Handle("setDelivered"))
	r.router.GET("/users/:user_uuid/conversations/:conv_id", r.Handle("getConversation"))
	r.router.POST("/users/:user_uuid/conversations/:conv_id/messages", r.Handle("sendMessage"))
	r.router.PUT("/users/:user_uuid/conversations/:conv_id/messages", r.Handle("setSeen"))
	r.router.GET("/users/:user_uuid/conversations/:conv_id/messages", r.Handle("getConversationMessages"))
	r.router.GET("/users/:user_uuid/conversations/:conv_id/messages/:mess_id", r.Handle("messageDetail"))
	r.router.DELETE("/users/:user_uuid/conversations/:conv_id/messages/:mess_id", r.Handle("deleteMessage"))
	r.router.GET("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", r.Handle("getMessageComments"))
	r.router.PUT("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", r.Handle("commentMessage"))
	r.router.DELETE("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", r.Handle("uncommentMessage"))
	r.router.POST("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/forward", r.Handle("forwardMessage"))

	return r.router
}

func (r *_router) Handle(endpointId string) httprouter.Handle {
	return r.wrap(r.app.GetEndpointHandler(endpointId))
}
