package api

import (
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	file, header, err := r.FormFile("profile_picture")
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	defer file.Close()

	filepath, err := SaveUploadedFile(file, header, "webui/public/images/profile", userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	var user database.User
	user.ID = userId
	user.ProfilePhotoURL = filepath
	result, err := h.db.UpdateUserProfile(&user)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	res := map[string]bool{"success": result}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
