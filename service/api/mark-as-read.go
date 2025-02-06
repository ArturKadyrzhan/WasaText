package api

import (
	"WasaText/service/database/models"
	"WasaText/service/helpers"
	"encoding/json"
	"net/http"
)

func (h *Handler) markAsRead(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)
	var input helpers.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	var conv *models.Conversation
	var err error

	if !input.IsGroup {
		conv, err = h.Repository.CheckPrivateConversation(userId, input.ToUserId)
		if err != nil {
			helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
			return
		}
	} else {
		conv, err = h.Repository.CheckGroupConversation(input.GroupId)
		if err != nil {
			helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
			return
		}
	}

	result, err := h.Repository.MarkAsRead(conv.ID, userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
