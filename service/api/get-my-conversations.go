package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := r.Context().Value("userId").(uint)
	info := make(map[string]interface{})
	users, err := h.Repository.GetConversationsUsers(userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	groups, err := h.Repository.GetGroups(userId)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
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
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

}
