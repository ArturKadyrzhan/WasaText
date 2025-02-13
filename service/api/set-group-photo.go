package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (h *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := r.FormValue("id")
	if id == "" {
		HandleError(w, NewAPIError("Group name is required", http.StatusBadRequest))
		return
	}

	file, header, err := r.FormFile("groupPhoto")
	if err != nil {
		HandleError(w, NewAPIError("Error retrieving file", http.StatusBadRequest))
		return
	}
	ID, _ := strconv.Atoi(id)
	filepath, err := SaveUploadedFile(file, header, "webui/public/images/group", uint(ID))

	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	result, err := h.db.UpdateGroupPhoto(ID, filepath)
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
