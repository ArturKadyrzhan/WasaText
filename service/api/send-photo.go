package api

import (
	"WasaText/service/helpers"
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) sendPhoto(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(uint)
	if !ok {
		helpers.HandleError(w, helpers.NewAPIError("Invalid user ID", http.StatusUnauthorized))
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Unable to parse form data", http.StatusUnprocessableEntity))
		return
	}

	isGroup, err := strconv.ParseBool(r.FormValue("isGroup"))
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Invalid 'isGroup' value", http.StatusUnprocessableEntity))
		return
	}

	groupId, err := parseOptionalInt(r.FormValue("groupId"))
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Invalid 'groupId' value", http.StatusUnprocessableEntity))
		return
	}

	toUserId, err := parseOptionalInt(r.FormValue("toUserId"))
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Invalid 'toUserId' value", http.StatusUnprocessableEntity))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Error retrieving file", http.StatusBadRequest))
		return
	}

	filepath, err := helpers.SaveUploadedFile(file, header, "webui/public/images/conversation", userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Error saving file", http.StatusBadRequest))
		return
	}

	payload := helpers.SendMessageRequest{
		ToUserId:  uint(toUserId),
		IsGroup:   isGroup,
		GroupId:   uint(groupId),
		PhotoPath: filepath,
	}

	result, err := h.SendMessage(userId, &payload)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Failed to send message", http.StatusBadRequest))
		return
	}

	messageResponse, err := helpers.SendMessageResponseHandler(&payload, userId, result, false)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	messageResponse["photoPath"] = filepath

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messageResponse)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
