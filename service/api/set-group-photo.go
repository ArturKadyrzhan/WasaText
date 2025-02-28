package api

import (
	"encoding/json"                       //
	"github.com/julienschmidt/httprouter" // library for handling HTTP routes.
	"net/http"                            // library for handling HTTP requests and responses.
	"strconv"                             //  Used to convert strings to numbers
)

func (h *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Extracting the Group ID from Request Data
	id := r.FormValue("id")
	if id == "" {
		HandleError(w, NewAPIError("Group name is required", http.StatusBadRequest))
		return
	}

	// Retrieving the Uploaded File
	file, header, err := r.FormFile("groupPhoto")
	if err != nil {
		HandleError(w, NewAPIError("Error retrieving file", http.StatusBadRequest))
		return
	}
	// convert,save
	ID, _ := strconv.Atoi(id)
	filepath, err := SaveUploadedFile(file, header, "webui/public/images/group", uint(ID))

	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	//update
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
