package api

import (
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input database.AddUsersToGroup
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}
	users, err := h.db.GetGroupMembers(input.GroupId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	var userIds = make([]uint, len(*users))

	for _, user := range *users {
		userIds = append(userIds, user.UserID)
	}

	for _, v := range input.Users {
		if !Contains(userIds, v.ID) {
			_, err = h.db.CreateGroupMembers(v.ID, userId, input.GroupId)
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

func Contains(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
