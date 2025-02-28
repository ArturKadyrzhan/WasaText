package api

import (
	"WasaText/service/database"           //   module for database operations.
	"encoding/json"                       //     Used to convert data into JSON
	"github.com/julienschmidt/httprouter" // library for handling HTTP routes.
	"net/http"                            // for handling HTTP requests and responses
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint) //  Converts the extracted value into an unsigned integer
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	// Extracting Comment Data from Request Body
	var input database.CommentMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	// Adding Comment to the Database
	result, err := h.db.CommentMessage(&input, userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	// formatting response then convert it to json
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
