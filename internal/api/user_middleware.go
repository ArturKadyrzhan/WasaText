package api

import (
	"WasaText/internal/helpers"
	"context"
	"log"
	"net/http"
	"strings"
)

func parseUserTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		log.Printf("Received Authorization Header: %s", token)

		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if token == "" {
			log.Printf("Authorization header missing")
			helpers.HandleError(w, helpers.NewAPIError("Authorization token is missing", http.StatusUnauthorized))
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		userId, err := helpers.ParseUserToken(token)

		if err != nil {
			log.Printf("Token validation error: %v", err)
			helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnauthorized))
			return
		}

		log.Printf("User ID from token: %d", userId)
		ctx := context.WithValue(r.Context(), "userId", userId)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
