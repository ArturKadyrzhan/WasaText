package api

import (
	"WasaText/service/database/models"
	"WasaText/service/helpers"
	"encoding/json"
	"net/http"
)

func (h *Handler) setMyPhoto(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	file, header, err := r.FormFile("profile_picture")
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	defer file.Close()

	filepath, err := helpers.SaveUploadedFile(file, header, "webui/public/images/profile", userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	var user models.User
	user.ID = userId
	user.ProfilePhotoURL = filepath
	result, err := h.Repository.UpdateUserProfile(&user)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	res := map[string]bool{"success": result}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
