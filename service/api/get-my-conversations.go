package api

import (
	"WasaText/service/helpers"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) getMyConversations(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)
	info := make(map[string]interface{})
	users, err := h.Repository.GetConversationsUsers(userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	groups, err := h.Repository.GetGroups(userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	fmt.Println(groups, "!!!! ERROR")
	info["users"] = users
	info["groups"] = groups

	res := map[string]interface{}{
		"result": info,
	}
	//
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

}
