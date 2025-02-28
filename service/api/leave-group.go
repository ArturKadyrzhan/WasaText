package api

import (
	"WasaText/service/database"           //  module for database operations.
	"encoding/json"                       //   Used to convert data into JSON.
	"github.com/julienschmidt/httprouter" // library for handling HTTP routes.
	"net/http"                            // for handling HTTP requests and responses
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var input database.GroupRequest
	// Defines a struct (input) to store the incoming JSON request body then extracts group id
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	userVal := r.Context().Value(keyUserID) // Extract the user ID stored in the request context.
	userId, ok := userVal.(uint)            //  Converts the extracted value into an unsigned integer
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	// Removing the User from the Group in the Database
	result, err := h.db.DeleteGroupMember(userId, input.Group)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	// Formatting the Response then convert to json
	res := map[string]bool{"success": result}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
