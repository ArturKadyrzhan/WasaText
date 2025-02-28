package api

import (
	"WasaText/service/database"           //  module for database operations.
	"encoding/json"                       //   Used to convert data into JSON.
	"github.com/julienschmidt/httprouter" //  library for handling HTTP routes.
	"net/http"                            // for handling HTTP requests and responses
)

func (h *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extracting Message Data from Request Body
	var input database.DeleteMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	// Deleting the Message from the Database
	result, err := h.db.DeleteMessage(input.MessageId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	// converting to json
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
