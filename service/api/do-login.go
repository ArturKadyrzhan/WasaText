package api

import (
	"WasaText/service/database/models"
	"WasaText/service/helpers"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) doLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	//fmt.Println(input, "error zdes' tt")
	existUser, err := h.Repository.GetUser(&user)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Failed to retrieve token", http.StatusBadRequest))
		return
	}
	var createdUser *models.User
	if existUser.ID == 0 {
		createdUser, err = h.Repository.CreateUser(&user)
		if err != nil {
			helpers.HandleError(w, helpers.NewAPIError("Failed to retrieve token", http.StatusBadRequest))
			return
		}

	}

	token, err := helpers.GenerateSessionToken(createdUser)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Failed to retrieve token", http.StatusBadRequest))
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
		helpers.HandleError(w, helpers.NewAPIError("Failed to encode response", http.StatusBadRequest))
		return
	}
}
