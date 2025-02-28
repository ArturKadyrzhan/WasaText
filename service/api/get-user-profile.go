package api

import (
	"WasaText/service/database"           // module for database operations.
	"encoding/json"                       // Used to convert data into JSON.
	"github.com/julienschmidt/httprouter" //  library for handling HTTP routes.
	"net/http"                            // library for handling HTTP requests and responses
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Extracting the User ID from the Request Context and error handling
	userVal := r.Context().Value(keyUserID) // Extracts the user ID stored in the request context.
	userId, ok := userVal.(uint)            // Converts the extracted value into an unsigned integer then handle error
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	// getting the search query from url
	profile, err := h.db.GetUserById(&database.User{ID: userId})
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	// mapping result and convert it to json
	res := map[string]database.User{"profile": *profile}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

}
