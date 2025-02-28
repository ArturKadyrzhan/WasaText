package api

import (
	"WasaText/service/database"           // module for database
	"encoding/json"                       // to convert data to json
	"github.com/julienschmidt/httprouter" // library for handling http routes
	"net/http"                            // for handling http requests and responses
)

func (h *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Extracting the User ID from the Request Context and error handling
	userVal := r.Context().Value(keyUserID) //  Extracts the user ID stored in the request context.
	userId, ok := userVal.(uint)            // Converts the extracted value into an unsigned integer
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Parsing Multipart Form Data (limit 10 MB)
	if err != nil {
		HandleError(w, NewAPIError("Unable to parse form data", http.StatusUnprocessableEntity))
		return
	}

	// Extracting group-name from form data
	groupName := r.FormValue("groupName")
	if groupName == "" {
		HandleError(w, NewAPIError("Group name is required", http.StatusBadRequest))
		return
	}

	// Extracting selected users from form data
	selectedUsers := r.FormValue("selectedUsers")
	if selectedUsers == "" {
		HandleError(w, NewAPIError("Selected users are required", http.StatusBadRequest))
		return
	}

	// Parsing the Selected Users List
	var users []database.User
	err = json.Unmarshal([]byte(selectedUsers), &users) // Converts the JSON string into a Go struct
	if err != nil {
		HandleError(w, NewAPIError("Invalid selected users data", http.StatusBadRequest))
		return
	}

	// Handling File Upload for Group Photo
	file, header, err := r.FormFile("profilePhoto")
	if err != nil {
		HandleError(w, NewAPIError("Error retrieving file", http.StatusBadRequest))
		return
	}
	filepath, err := SaveUploadedFile(file, header, "webui/public/images/group", userId)

	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	// Creating the Group in the Database
	payload := &database.CreateGroupRequest{
		GroupName:      groupName,
		GroupPhotoPath: filepath,
		Users:          users,
	}
	group, err := h.db.CreateGroup(payload, userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	_, err = h.db.CreateGroupMembers(userId, userId, group.ID)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	// Adding Members to the Group
	for _, user := range payload.Users {
		_, err = h.db.CreateGroupMembers(user.ID, userId, group.ID)
		if err != nil {
			HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
			return
		}
	}

	res := map[string]bool{"success": true}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
