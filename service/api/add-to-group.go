package api

import (
	"WasaText/service/helpers"
	"encoding/json"
	"net/http"
	"slices"
)

func (h *Handler) addToGroup(w http.ResponseWriter, r *http.Request) {
	var input helpers.AddUsersToGroup
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	userId := r.Context().Value("userId").(uint)

	users, err := h.Repository.GetGroupMembers(input.GroupId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	var userIds = make([]uint, len(*users))

	for _, user := range *users {
		userIds = append(userIds, user.UserID)
	}

	for _, v := range input.Users {
		if !slices.Contains(userIds, v.ID) {
			_, err = h.Repository.CreateGroupMembers(v.ID, userId, input.GroupId)
			if err != nil {
				helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
				return
			}
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
