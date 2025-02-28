package api

import (
	"WasaText/service/api/reqcontext"     // Context for handling request details
	"WasaText/service/database"           // Database operations
	"encoding/json"                       // to convert it json
	"github.com/julienschmidt/httprouter" /// Library for handling HTTP routes
	"net/http"                            // package for HTTP requests and responses
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extracting User Credentials from Request Body
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	// Checking If User Exists in Database
	existUser, err := h.db.GetUser(&user)
	if err != nil {
		HandleError(w, NewAPIError("Failed to retrieve token", http.StatusBadRequest))
		return
	}

	// Creating a New User If Not Found
	if existUser.ID == 0 {
		existUser, err = h.db.CreateUser(&user)
		if err != nil {
			HandleError(w, NewAPIError("Failed to retrieve token", http.StatusBadRequest))
			return
		}

	}

	// Generating an Authentication Token
	token, err := GenerateSessionToken(existUser)
	if err != nil {
		HandleError(w, NewAPIError("Failed to retrieve token", http.StatusBadRequest))
		return
	}

	// Formatting the Response
	res := map[string]interface{}{
		"token":    token,
		"id":       user.ID,
		"username": user.Username,
	}

	// Sending the Response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError("Failed to encode response", http.StatusBadRequest))
		return
	}
}
