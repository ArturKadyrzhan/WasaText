package api

import (
	"WasaText/service/database" // database operations
	"encoding/json"
	"github.com/julienschmidt/httprouter" //  library for handling HTTP routes.
	"net/http"                            // for handling HTTP requests and responses
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,

func (h *_router) markAsRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extracting the User ID from Request Context
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	// Extracting Conversation structure from Request Body
	var input database.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	// Checking if the Conversation Exists
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
	// Marking Messages as Read in the Database
	result, err := h.db.MarkAsRead(conv.ID, userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	// Formatting the Response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
