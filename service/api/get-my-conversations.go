package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	info := make(map[string]interface{})
	users, err := h.db.GetConversationsUsers(userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	groups, err := h.db.GetGroups(userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	info["users"] = users
	info["groups"] = groups

	res := map[string]interface{}{
		"result": info,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

}
