package api

import (
	"WasaText/service/helpers"
	"encoding/json"
	"net/http"
)

func (h *Handler) deleteMessage(w http.ResponseWriter, r *http.Request) {
	var input helpers.DeleteMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	result, err := h.Repository.DeleteMessage(input.MessageId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
