package api

import (
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}
	query := r.URL.Query().Get("search")
	if query == "" {
		HandleError(w, NewAPIError("Query is empty", http.StatusBadRequest))
		return
	}

	users, err := h.db.GetUsers(query, userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	res := map[string][]database.User{"users": *users}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

}
