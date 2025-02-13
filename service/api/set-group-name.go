package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	groupName := struct {
		ID        int    `json:"id"`
		GroupName string `json:"groupName"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&groupName); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	result, err := h.db.UpdateGroupName(groupName.ID, groupName.GroupName)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
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
