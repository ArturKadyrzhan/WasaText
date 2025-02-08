package api

import (
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"slices"
)

func (h *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input database.AddUsersToGroup
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	userId := r.Context().Value("userId").(uint)

	users, err := h.Repository.GetGroupMembers(input.GroupId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
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
				HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
				return
			}
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
