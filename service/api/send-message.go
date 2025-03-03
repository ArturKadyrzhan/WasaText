package api

import (
	"WasaText/service/database"           //  module for database operations.
	"encoding/json"                       //  Used to convert data into JSON.
	"github.com/julienschmidt/httprouter" // library for handling HTTP routes.
	"net/http"                            //   for handling HTTP requests and responses
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Extracting Message Data from Request Body
	var input database.SendMessageRequest // Defines a struct (input) to store the incoming JSON request body.
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	// Extracting the Sender's User ID
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint) //Converts the extracted value into an unsigned integer
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	//Processing the Message
	result, err := h.SendMessage(userId, &input)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	// Handling the Response for the Sent Message
	messageResponse, err := SendMessageResponseHandler(&input, userId, result, false)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messageResponse)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
