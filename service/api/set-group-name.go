package api

import (
	"encoding/json"                       // to convert our data to json
	"github.com/julienschmidt/httprouter" // library for handling http routes
	"net/http"                            // for handling http requests and responses
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,

func (h *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Defines a struct (groupName) to store the incoming JSON request body.
	groupName := struct {
		ID        int    `json:"id"`
		GroupName string `json:"groupName"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&groupName); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	// Updating the Group Name in Database
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
