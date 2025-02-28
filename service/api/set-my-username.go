package api

import (
	"WasaText/service/database"           // module for database
	"encoding/json"                       // converting our data to json
	"github.com/julienschmidt/httprouter" // library to handle http routes
	"net/http"                            //  library for handling HTTP requests and responses.
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userVal := r.Context().Value(keyUserID) // Converts the extracted value into an unsigned integer then handle error
	userId, ok := userVal.(uint)            // Converts the extracted value into an unsigned integer then handling error
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}
	// Extracting the New Username from Request Body and handle error
	username := struct {
		Username string `json:"username"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	// Fetching the User from the Database
	var user database.User
	user.ID = userId
	userProfile, err := h.db.GetUserById(&user)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	// Updating the Username
	userProfile.Username = username.Username
	_, err = h.db.UpdateUserProfile(&user)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	// Formatting the Response as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(username)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
