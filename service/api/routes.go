package api

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/health/alive", logRequests(http.HandlerFunc(healthcheck)))
	mux.Handle("/health/ready", logRequests(http.HandlerFunc(healthcheck)))
	mux.Handle("/login", logRequests(handleRequest(http.MethodPost, h.doLogin)))
	mux.Handle("/users", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodGet, h.getUsers))))
	mux.Handle("/conversations", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodGet, h.getMyConversations))))
	mux.Handle("/profile/photo", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.setMyPhoto))))
	mux.Handle("/group", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.createGroup))))
	mux.Handle("/group/users", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.addToGroup))))
	mux.Handle("/group/leave", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.leaveGroup))))
	mux.Handle("/message", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.sendMessage))))
	mux.Handle("/message", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodDelete, h.deleteMessage))))
	mux.Handle("/messages", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodGet, h.getMessages))))
	mux.Handle("/message/read", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.markAsRead))))
	mux.Handle("/message/send-photo", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.sendPhoto))))
	mux.Handle("/message/comment", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.commentMessage))))
	mux.Handle("/message/uncomment", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.uncommentMessage))))
	mux.Handle("/message/forward", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.forwardMessage))))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})

}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"status":  http.StatusOK,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleRequest(method string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			handlerFunc(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
