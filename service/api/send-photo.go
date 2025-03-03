package api

import (
	"WasaText/service/database"           //  module for database operations.
	"encoding/json"                       //    Used to convert data into JSON.
	"github.com/julienschmidt/httprouter" // library for handling HTTP routes.
	"net/http"                            // for handling HTTP requests and responses
	"strconv"                             //  for converting strings to int/bool
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) sendPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Extracting the User ID from Request Context
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint) // Converts the extracted value into an unsigned integer
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	// Parsing Multipart Form Data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		HandleError(w, NewAPIError("Unable to parse form data", http.StatusUnprocessableEntity))
		return
	}

	// Extracting Conversation Details
	isGroup, err := strconv.ParseBool(r.FormValue("isGroup"))
	if err != nil {
		HandleError(w, NewAPIError("Invalid 'isGroup' value", http.StatusUnprocessableEntity))
		return
	}

	groupId, err := parseOptionalInt(r.FormValue("groupId"))
	if err != nil {
		HandleError(w, NewAPIError("Invalid 'groupId' value", http.StatusUnprocessableEntity))
		return
	}

	toUserId, err := parseOptionalInt(r.FormValue("toUserId"))
	if err != nil {
		HandleError(w, NewAPIError("Invalid 'toUserId' value", http.StatusUnprocessableEntity))
		return
	}

	// Retrieving the Uploaded File
	file, header, err := r.FormFile("file")
	if err != nil {
		HandleError(w, NewAPIError("Error retrieving file", http.StatusBadRequest))
		return
	}

	// Saving the Uploaded File
	filepath, err := SaveUploadedFile(file, header, "webui/public/images/conversation", userId)
	if err != nil {
		HandleError(w, NewAPIError("Error saving file", http.StatusBadRequest))
		return
	}

	// prepare message structure before sending
	payload := database.SendMessageRequest{
		ToUserId:  uint(toUserId),
		IsGroup:   isGroup,
		GroupId:   uint(groupId),
		PhotoPath: filepath,
	}

	// Sending  Photo Message
	result, err := h.SendMessage(userId, &payload) // Stores the photo message in the database.
	if err != nil {
		HandleError(w, NewAPIError("Failed to send message", http.StatusBadRequest))
		return
	}

	// Formatting the Response
	messageResponse, err := SendMessageResponseHandler(&payload, userId, result, false)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	messageResponse["photoPath"] = filepath //Adds photoPath to the response.

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messageResponse)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
