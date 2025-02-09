package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.Handle("doLogin"))
	rt.router.GET("/users", rt.Handle("getUsers"))
	rt.router.GET("/users/:user_uuid", rt.Handle("getUser"))
	rt.router.PUT("/users/:user_uuid/username", rt.Handle("setMyUserName"))
	rt.router.PUT("/users/:user_uuid/avatar", rt.Handle("setMyPhoto"))
	rt.router.POST("/users/:user_uuid/chats", rt.Handle("createChat"))
	rt.router.POST("/users/:user_uuid/groups", rt.Handle("createGroup"))
	rt.router.PUT("/users/:user_uuid/groups/:conv_id", rt.Handle("addToGroup"))
	rt.router.DELETE("/users/:user_uuid/groups/:conv_id", rt.Handle("leaveGroup"))
	rt.router.PUT("/users/:user_uuid/groups/:conv_id/name", rt.Handle("setGroupName"))
	rt.router.PUT("/users/:user_uuid/groups/:conv_id/photo", rt.Handle("setGroupPhoto"))
	rt.router.GET("/users/:user_uuid/conversations", rt.Handle("getMyConversations"))
	rt.router.PUT("/users/:user_uuid/conversations", rt.Handle("setDelivered"))
	rt.router.GET("/users/:user_uuid/conversations/:conv_id", rt.Handle("getConversation"))
	rt.router.POST("/users/:user_uuid/conversations/:conv_id/messages", rt.Handle("sendMessage"))
	rt.router.PUT("/users/:user_uuid/conversations/:conv_id/messages", rt.Handle("setSeen"))
	rt.router.GET("/users/:user_uuid/conversations/:conv_id/messages", rt.Handle("getConversationMessages"))
	rt.router.GET("/users/:user_uuid/conversations/:conv_id/messages/:mess_id", rt.Handle("messageDetail"))
	rt.router.DELETE("/users/:user_uuid/conversations/:conv_id/messages/:mess_id", rt.Handle("deleteMessage"))
	rt.router.GET("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", rt.Handle("getMessageComments"))
	rt.router.PUT("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", rt.Handle("commentMessage"))
	rt.router.DELETE("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/comments", rt.Handle("uncommentMessage"))
	rt.router.POST("/users/:user_uuid/conversations/:conv_id/messages/:mess_id/forward", rt.Handle("forwardMessage"))

	return rt.router
}

func (rt *_router) Handle(endpointId string) httprouter.Handle {
	return rt.wrap(rt.app.GetEndpointHandler(endpointId))
}
