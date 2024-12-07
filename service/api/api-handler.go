package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/session", rt.ResponseHeaderMiddleware(rt.doLogin))
	rt.router.GET("/users", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.getUsers)))
	rt.router.GET("/users/:username", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.getUser)))
	rt.router.PUT("/users/:username/username", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.setMyUsername)))
	rt.router.PUT("/users/:username/photo", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.setMyPhoto)))
	//
	rt.router.GET("/users/:username/conversations", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.getMyConversations)))
	//rt.router.POST("/users/:username/groups", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.createGroup)))
	//rt.router.PUT("/users/:username/groups/:group", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.addToGroup)))
	//rt.router.DELETE("/users/:username/groups/:group", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.leaveGroup)))
	//rt.router.PUT("/users/:username/groups/:group/name", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.setGroupName)))
	//rt.router.PUT("/users/:username/groups/:group/photo", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.setGroupPhoto)))
	//
	//rt.router.POST("/users/:username/chats", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.createChat)))
	//rt.router.GET("/users/:username/test", rt.ResponseHeaderMiddleware(rt.AuthenticationMiddleware(rt.test)))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
