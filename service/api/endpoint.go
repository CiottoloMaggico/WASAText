package api

import (
	"github.com/ciottolomaggico/wasatext/service/api/routes"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetSessionRouter()
	return router.DoLogin(w, r, params, context)
}

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetUserRouter()
	return router.GetUsers(w, r, params, context)
}

func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetUserRouter()
	return router.GetUser(w, r, params, context)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetUserRouter()
	return router.SetMyUsername(w, r, params, context)
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetUserRouter()
	return router.SetMyPhoto(w, r, params, context)
}

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetConversationRouter()
	return router.CreateChat(w, r, params, context)
}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetConversationRouter()
	return router.CreateGroup(w, r, params, context)
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetConversationRouter()
	return router.AddToGroup(w, r, params, context)
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetConversationRouter()
	return router.LeaveGroup(w, r, params, context)
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetConversationRouter()
	return router.SetGroupName(w, r, params, context)
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetConversationRouter()
	return router.SetGroupPhoto(w, r, params, context)
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetUserConversationRouter()
	return router.GetMyConversations(w, r, params, context)
}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetUserConversationRouter()
	return router.GetConversation(w, r, params, context)
}

func (rt *_router) setDelivered(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetMessageRouter()
	return router.SetDelivered(w, r, params, context)
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetMessageRouter()
	return router.SendMessage(w, r, params, context)
}

func (rt *_router) setSeen(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetMessageRouter()
	return router.SetSeen(w, r, params, context)
}

func (rt *_router) getConversationMessages(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetMessageRouter()
	return router.GetConversationMessages(w, r, params, context)
}

func (rt *_router) messageDetail(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetMessageRouter()
	return router.GetConversationMessageDetail(w, r, params, context)
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetMessageRouter()
	return router.DeleteConversationMessage(w, r, params, context)
}

func (rt *_router) getMessageComments(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetMessageInfoRouter()
	return router.GetMessageComments(w, r, params, context)
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetMessageInfoRouter()
	return router.SetMessageComment(w, r, params, context)
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetMessageInfoRouter()
	return router.RemoveMessageComment(w, r, params, context)
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context routes.RequestContext) error {
	router := rt.app.GetMessageRouter()
	return router.ForwardMessage(w, r, params, context)
}
