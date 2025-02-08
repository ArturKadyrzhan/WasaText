package api

import (
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) markAsRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := r.Context().Value("userId").(uint)
	var input database.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	var conv *database.Conversation
	var err error

	if !input.IsGroup {
		conv, err = h.Repository.CheckPrivateConversation(userId, input.ToUserId)
		if err != nil {
			HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
			return
		}
	} else {
		conv, err = h.Repository.CheckGroupConversation(input.GroupId)
		if err != nil {
			HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
			return
		}
	}

	result, err := h.Repository.MarkAsRead(conv.ID, userId)
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
