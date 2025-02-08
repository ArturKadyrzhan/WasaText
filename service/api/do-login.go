package api

import (
	"WasaText/service/database"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	//fmt.Println(input, "error zdes' tt")
	existUser, err := h.Repository.GetUser(&user)
	if err != nil {
		HandleError(w, NewAPIError("Failed to retrieve token", http.StatusBadRequest))
		return
	}
	var createdUser *database.User
	if existUser.ID == 0 {
		createdUser, err = h.Repository.CreateUser(&user)
		if err != nil {
			HandleError(w, NewAPIError("Failed to retrieve token", http.StatusBadRequest))
			return
		}

	}

	token, err := GenerateSessionToken(createdUser)
	if err != nil {
		HandleError(w, NewAPIError("Failed to retrieve token", http.StatusBadRequest))
		return
	}

	res := map[string]interface{}{
		"token":    token,
		"id":       user.ID,
		"username": user.Username,
	}

	fmt.Println(token)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError("Failed to encode response", http.StatusBadRequest))
		return
	}
}
