package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	rt.router.GET("/liveness", rt.healthcheck)
	rt.router.POST("/login", rt.doLogin)
	rt.router.GET("/users", rt.parseUserTokenMiddleware(rt.getUsers))
	rt.router.GET("/conversations", rt.getMyConversations)
	rt.router.POST("/profile/photo", rt.setMyPhoto)
	rt.router.POST("/group", rt.createGroup)
	rt.router.POST("/group/users", rt.addToGroup)
	rt.router.POST("/group/leave", rt.leaveGroup)
	rt.router.POST("/message", rt.sendMessage)
	rt.router.DELETE("/message", rt.deleteMessage)
	rt.router.GET("/messages", rt.getMessages)
	rt.router.POST("/message/read", rt.markAsRead)
	rt.router.POST("/message/send-photo", rt.sendPhoto)
	rt.router.POST("/message/comment", rt.commentMessage)
	rt.router.POST("/message/uncomment", rt.uncommentMessage)
	rt.router.POST("/message/forward", rt.forwardMessage)
	return rt.router

}
