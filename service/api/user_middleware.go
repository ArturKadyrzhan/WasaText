package api

import (
	"context"                             // for request context
	"github.com/julienschmidt/httprouter" // for handling http routes
	"log"                                 // for log events
	"net/http"
	"strings" // for string manipulation functions
)

// Defining a Context Key for User ID
type contextKey string

const keyUserID contextKey = "userId"

// Defines a new type (httpMiddlewareHandler) for functions that handle HTTP requests.
type httpMiddlewareHandler func(http.ResponseWriter, *http.Request, httprouter.Params)

// Defining the Middleware Function
func (rt *_router) parseUserTokenMiddleware(fn httpMiddlewareHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Extracting the Token from Request Headers
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

		// Parsing the Token to Extract User ID
		userId, err := ParseUserToken(token)
		if err != nil {
			HandleError(w, NewAPIError(err.Error(), http.StatusUnauthorized))
			return
		}

		// Storing the User ID in the Request Context
		log.Printf("User ID from token: %d", userId)

		ctx := context.WithValue(r.Context(), keyUserID, userId)
		r = r.WithContext(ctx)

		fn(w, r, ps)
	}
}
