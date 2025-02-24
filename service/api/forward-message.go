package api

import (
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	var input database.ForwardMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	if input.IsPhoto {
		for _, user := range input.Users {
			sendMessagePayload := database.SendMessageRequest{
				ToUserId:      user.ID,
				IsGroup:       false,
				PhotoPath:     input.Text,
				ForwardedById: input.ForwardedById,
			}
			_, err := h.SendMessage(userId, &sendMessagePayload)
			if err != nil {
				HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
				return
			}
		}
	} else {
		for _, user := range input.Users {
			sendMessagePayload := database.SendMessageRequest{
				ToUserId:      user.ID,
				IsGroup:       false,
				Text:          input.Text,
				ForwardedById: input.ForwardedById,
			}

			_, err := h.SendMessage(userId, &sendMessagePayload)
			if err != nil {
				HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(true)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
