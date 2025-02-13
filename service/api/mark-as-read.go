package api

import (
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) markAsRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	var input database.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	var conv *database.Conversation
	var err error

	if !input.IsGroup {
		conv, err = h.db.CheckPrivateConversation(userId, input.ToUserId)
		if err != nil {
			HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
			return
		}
	} else {
		conv, err = h.db.CheckGroupConversation(input.GroupId)
		if err != nil {
			HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
			return
		}
	}

	result, err := h.db.MarkAsRead(conv.ID, userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
