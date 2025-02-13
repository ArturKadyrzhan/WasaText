package api

import (
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}
	username := struct {
		Username string `json:"username"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	var user database.User
	user.ID = userId
	userProfile, err := h.db.GetUserById(&user)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	userProfile.Username = username.Username
	_, err = h.db.UpdateUserProfile(&user)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(username)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
