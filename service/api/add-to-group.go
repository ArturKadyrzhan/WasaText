package api

import (
	"WasaText/service/database"           //  module for database operations.
	"encoding/json"                       // Used to convert data into JSON.
	"github.com/julienschmidt/httprouter" // library for handling HTTP routes.
	"net/http"                            // for handling HTTP requests and responses
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var input database.AddUsersToGroup
	// Defines a struct (input) to store the incoming JSON request body.
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	// Extracting the group ID and users to be added from the request.
	userVal := r.Context().Value(keyUserID) // Extracts the user ID stored in the request context.
	userId, ok := userVal.(uint)            // Converts the extracted value into an unsigned integer then handle error
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}
	// getting Fetching Existing Group Members
	users, err := h.db.GetGroupMembers(input.GroupId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	// Creating a List of Current Group Member IDs
	var userIds = make([]uint, len(*users))
	for _, user := range *users {
		userIds = append(userIds, user.UserID)
	}

	// Adding New Users to the Group
	for _, v := range input.Users {
		if !Contains(userIds, v.ID) {
			_, err = h.db.CreateGroupMembers(v.ID, userId, input.GroupId)
			if err != nil {
				HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
				return
			}
		}
	}

	// formatting response then convert to json
	res := map[string]bool{"success": true}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}

// Helper Function: Checking if a User is Already in the Group
func Contains(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
