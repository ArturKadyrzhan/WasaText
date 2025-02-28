package api

import (
	"WasaText/service/database"
	"encoding/json"                       // to convert our data to json
	"github.com/julienschmidt/httprouter" // library for handling http routes
	"net/http"                            // for handling http requests and responses
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extracting the User ID from the Request Context and error handling
	userVal := r.Context().Value(keyUserID)
	// Converts the extracted value into an unsigned integer then handle error
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	// Retrieving the Uploaded File
	file, header, err := r.FormFile("profile_picture")
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	// Ensures the file is closed after processing.
	defer file.Close()

	// Save uploaded file and handling error
	filepath, err := SaveUploadedFile(file, header, "webui/public/images/profile", userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	// Fetching the User from the Database then updating it and handling errors
	user := database.User{
		ID: userId,
	}
	updateUser, err := h.db.GetUserById(&user)
	user.ProfilePhotoURL = filepath
	user.Username = updateUser.Username
	_, err = h.db.UpdateUserProfile(&user)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	// Creating the Response Object then formatting response to json
	res := map[string]string{"success": filepath}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
