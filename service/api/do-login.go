package api

import (
	"WasaText/service/api/reqcontext"
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	existUser, err := h.db.GetUser(&user)
	if err != nil {
		HandleError(w, NewAPIError("Failed to retrieve token", http.StatusBadRequest))
		return
	}

	if existUser.ID == 0 {
		existUser, err = h.db.CreateUser(&user)
		if err != nil {
			HandleError(w, NewAPIError("Failed to retrieve token", http.StatusBadRequest))
			return
		}

	}

	token, err := GenerateSessionToken(existUser)
	if err != nil {
		HandleError(w, NewAPIError("Failed to retrieve token", http.StatusBadRequest))
		return
	}

	res := map[string]interface{}{
		"token":    token,
		"id":       user.ID,
		"username": user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError("Failed to encode response", http.StatusBadRequest))
		return
	}
}
