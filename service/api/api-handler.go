package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	rt.router.GET("/liveness", rt.healthcheck)
	rt.router.POST("/session", rt.doLogin)
	rt.router.GET("/users", rt.parseUserTokenMiddleware(rt.getUsers))
	rt.router.GET("/conversations", rt.parseUserTokenMiddleware(rt.getMyConversations))
	rt.router.POST("/profile/photo", rt.parseUserTokenMiddleware(rt.setMyPhoto))
	rt.router.POST("/group", rt.parseUserTokenMiddleware(rt.createGroup))
	rt.router.POST("/group/users", rt.parseUserTokenMiddleware(rt.addToGroup))
	rt.router.POST("/group/leave", rt.parseUserTokenMiddleware(rt.leaveGroup))
	rt.router.POST("/message/send", rt.parseUserTokenMiddleware(rt.sendMessage))
	rt.router.POST("/message", rt.parseUserTokenMiddleware(rt.deleteMessage))
	rt.router.POST("/messages", rt.parseUserTokenMiddleware(rt.getMessages))
	rt.router.POST("/message/read", rt.parseUserTokenMiddleware(rt.markAsRead))
	rt.router.POST("/message/send-photo", rt.parseUserTokenMiddleware(rt.sendPhoto))
	rt.router.POST("/message/comment", rt.parseUserTokenMiddleware(rt.commentMessage))
	rt.router.POST("/message/uncomment", rt.parseUserTokenMiddleware(rt.uncommentMessage))
	rt.router.POST("/message/forward", rt.parseUserTokenMiddleware(rt.forwardMessage))
	return rt.router

}
