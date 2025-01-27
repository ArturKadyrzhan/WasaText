package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/health/alive", logRequests(http.HandlerFunc(healthcheck)))
	mux.Handle("/health/ready", logRequests(http.HandlerFunc(healthcheck)))
	mux.Handle("/login", logRequests(handleRequest(http.MethodPost, h.doLogin)))
	mux.Handle("/get-users", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodGet, h.getUsers))))
	mux.Handle("/get-conversations", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.getConversations))))
	mux.Handle("/send-message", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.sendMessage))))
	mux.Handle("/get-messages", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.getMessages))))
	mux.Handle("/mark-as-read", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.markAsRead))))
	mux.Handle("/create-group", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.createGroup))))
	mux.Handle("/send-photo", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.sendPhoto))))
	mux.Handle("/upload-profile-picture", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.uploadProfilePicture))))
	mux.Handle("/update-group", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.addUsersToGroup))))
	mux.Handle("/delete-message", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.deleteMessage))))
	mux.Handle("/leave-group", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.leaveGroup))))
	mux.Handle("/comment-message", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.commentMessage))))
	mux.Handle("/uncomment-message", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.uncommentMessage))))
	mux.Handle("/forward-message", logRequests(parseUserTokenMiddleware(handleRequest(http.MethodPost, h.forwardMessage))))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})

}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"status":  http.StatusOK,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	log.Printf("CORS headers set for %s", r.Method)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "1")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Allowing requests from port 5173
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
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
