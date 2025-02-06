package api

import (
	"WasaText/service/helpers"
	"encoding/json"
	"net/http"
)

func (h *Handler) commentMessage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	var input helpers.CommentMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	result, err := h.Repository.CommentMessage(&input, userId)
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
