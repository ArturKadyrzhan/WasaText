package api

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const keyUserID contextKey = "userId"

type httpMiddlewareHandler func(http.ResponseWriter, *http.Request, httprouter.Params)

func (rt *_router) parseUserTokenMiddleware(fn httpMiddlewareHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("Authorization")
		log.Printf("Received Authorization Header: %s", token)

		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if token == "" {
			log.Printf("Authorization header missing")
			HandleError(w, NewAPIError("Authorization token is missing", http.StatusUnauthorized))
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		userId, err := ParseUserToken(token)
		if err != nil {
			HandleError(w, NewAPIError(err.Error(), http.StatusUnauthorized))
			return
		}

		log.Printf("User ID from token: %d", userId)

		ctx := context.WithValue(r.Context(), keyUserID, userId)
		r = r.WithContext(ctx)

		fn(w, r, ps)
	}
}
