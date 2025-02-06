package api

import (
	"WasaText/service/helpers"
	"encoding/json"
	"net/http"
)

func (h *Handler) forwardMessage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	var input helpers.ForwardMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	if input.IsPhoto {
		for _, user := range input.Users {
			sendMessagePayload := helpers.SendMessageRequest{
				ToUserId:  user.ID,
				IsGroup:   false,
				PhotoPath: input.Text,
			}
			_, err := h.SendMessage(userId, &sendMessagePayload)
			if err != nil {
				helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
				return
			}
		}
	} else {
		for _, user := range input.Users {
			sendMessagePayload := helpers.SendMessageRequest{
				ToUserId: user.ID,
				IsGroup:  false,
				Text:     input.Text,
			}

			_, err := h.SendMessage(userId, &sendMessagePayload)
			if err != nil {
				helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(true)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
