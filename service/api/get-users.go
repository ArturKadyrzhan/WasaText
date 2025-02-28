package api

import (
	"WasaText/service/database"
	"encoding/json"                       // – I use to convert data into JSON format.
	"github.com/julienschmidt/httprouter" // library to handle routes in HTTP requests.
	"net/http"                            // standard Go library for handling HTTP requests and responses.
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,

func (h *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extracting the User ID from the Request Context and error handling
	userVal := r.Context().Value(keyUserID) // Extracts the user ID stored in the request context.
	userId, ok := userVal.(uint)            // Converts the extracted value into an unsigned integer then handle error
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	// getting the search query from url
	query := r.URL.Query().Get("search")
	if query == "" {
		HandleError(w, NewAPIError("Query is empty", http.StatusBadRequest))
		return
	}

	// fetching users from database
	users, err := h.db.GetUsers(query, userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	// finally get result then format it as JSON
	res := map[string][]database.User{"users": *users}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

}
