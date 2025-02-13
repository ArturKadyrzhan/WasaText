package api

import (
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input database.GroupRequest
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

	result, err := h.db.DeleteGroupMember(userId, input.Group)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusUnprocessableEntity))
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
