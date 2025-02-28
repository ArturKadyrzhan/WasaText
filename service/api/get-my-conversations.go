package api

import (
	"encoding/json"                       // Used to convert data into JSON.
	"github.com/julienschmidt/httprouter" // for handling HTTP routes.
	"net/http"                            // for handling HTTP requests and responses.
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Extracting the User ID from the Request Context and error handling
	userVal := r.Context().Value(keyUserID) // Extracts the user ID stored in the request context.
	userId, ok := userVal.(uint)            // Converts the extracted value into an unsigned integer then handle error
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	// Create map to store conversations then fetch private conversations
	info := make(map[string]interface{})
	users, err := h.db.GetConversationsUsers(userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	// Also fetch group conversations
	groups, err := h.db.GetGroups(userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	// Storing conversations into the map then create res object
	info["users"] = users
	info["groups"] = groups

	res := map[string]interface{}{
		"result": info,
	}

	// finalLy formatting res to JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

}
