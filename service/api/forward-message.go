package api

import (
	"WasaText/service/database" // for database module operations
	"encoding/json"
	"github.com/julienschmidt/httprouter" // for handling http routes
	"net/http"                            // go package for http requests
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extracting the User ID from Request Context
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	// Extracting ForwardMessage Data from Request Body
	var input database.ForwardMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	// Forwarding a Photo Message
	if input.IsPhoto {
		for _, user := range input.Users {
			sendMessagePayload := database.SendMessageRequest{
				ToUserId:      user.ID,
				IsGroup:       false,
				PhotoPath:     input.Text,
				ForwardedById: int(userId),
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
				ForwardedById: int(userId),
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
