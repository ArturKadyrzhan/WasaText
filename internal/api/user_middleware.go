package api

import (
	"WasaText/internal/helpers"
	"context"
	"net/http"
	"strings"
)

func parseUserTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if token == "" {
			helpers.HandleError(w, helpers.NewAPIError("Authorization token is missing", http.StatusUnauthorized))
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		userId, err := helpers.ParseUserToken(token)
		if err != nil {
			helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnauthorized))
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
