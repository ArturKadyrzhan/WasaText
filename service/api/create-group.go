package api

import (
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func NewHandler(repository *database.Repository) *Handler {
	return &Handler{Repository: repository}
}

func (h *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := r.Context().Value("userId").(uint)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		HandleError(w, NewAPIError("Unable to parse form data", http.StatusUnprocessableEntity))
		return
	}

	groupName := r.FormValue("groupName")
	if groupName == "" {
		HandleError(w, NewAPIError("Group name is required", http.StatusBadRequest))
		return
	}

	selectedUsers := r.FormValue("selectedUsers")
	if selectedUsers == "" {
		HandleError(w, NewAPIError("Selected users are required", http.StatusBadRequest))
		return
	}

	var users []database.User
	err = json.Unmarshal([]byte(selectedUsers), &users)
	if err != nil {
		HandleError(w, NewAPIError("Invalid selected users data", http.StatusBadRequest))
		return
	}

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
	payload := &database.CreateGroupRequest{
		GroupName:      groupName,
		GroupPhotoPath: filepath,
		Users:          users,
	}
	group, err := h.Repository.CreateGroup(payload, userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	_, err = h.Repository.CreateGroupMembers(userId, userId, group.ID)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	for _, user := range payload.Users {
		_, err = h.Repository.CreateGroupMembers(user.ID, userId, group.ID)
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
