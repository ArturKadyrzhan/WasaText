package api

import (
	"WasaText/service/database"
	"WasaText/service/database/models"
	"WasaText/service/helpers"
	"encoding/json"
	"net/http"
)

func NewHandler(repository *database.Repository) *Handler {
	return &Handler{Repository: repository}
}

func (h *Handler) createGroup(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Unable to parse form data", http.StatusUnprocessableEntity))
		return
	}

	groupName := r.FormValue("groupName")
	if groupName == "" {
		helpers.HandleError(w, helpers.NewAPIError("Group name is required", http.StatusBadRequest))
		return
	}

	selectedUsers := r.FormValue("selectedUsers")
	if selectedUsers == "" {
		helpers.HandleError(w, helpers.NewAPIError("Selected users are required", http.StatusBadRequest))
		return
	}

	var users []models.User
	err = json.Unmarshal([]byte(selectedUsers), &users)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Invalid selected users data", http.StatusBadRequest))
		return
	}

	file, header, err := r.FormFile("profilePhoto")
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Error retrieving file", http.StatusBadRequest))
		return
	}
	filepath, err := helpers.SaveUploadedFile(file, header, "webui/public/images/group", userId)

	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	payload := &helpers.CreateGroupRequest{
		GroupName:      groupName,
		GroupPhotoPath: filepath,
		Users:          users,
	}
	group, err := h.Repository.CreateGroup(payload, userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	_, err = h.Repository.CreateGroupMembers(userId, userId, group.ID)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	for _, user := range payload.Users {
		_, err = h.Repository.CreateGroupMembers(user.ID, userId, group.ID)
		if err != nil {
			helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
			return
		}
	}

	res := map[string]bool{"success": true}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
