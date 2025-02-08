package api

import (
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (h *_router) sendPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, ok := r.Context().Value("userId").(uint)
	if !ok {
		HandleError(w, NewAPIError("Invalid user ID", http.StatusUnauthorized))
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		HandleError(w, NewAPIError("Unable to parse form data", http.StatusUnprocessableEntity))
		return
	}

	isGroup, err := strconv.ParseBool(r.FormValue("isGroup"))
	if err != nil {
		HandleError(w, NewAPIError("Invalid 'isGroup' value", http.StatusUnprocessableEntity))
		return
	}

	groupId, err := parseOptionalInt(r.FormValue("groupId"))
	if err != nil {
		HandleError(w, NewAPIError("Invalid 'groupId' value", http.StatusUnprocessableEntity))
		return
	}

	toUserId, err := parseOptionalInt(r.FormValue("toUserId"))
	if err != nil {
		HandleError(w, NewAPIError("Invalid 'toUserId' value", http.StatusUnprocessableEntity))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		HandleError(w, NewAPIError("Error retrieving file", http.StatusBadRequest))
		return
	}

	filepath, err := SaveUploadedFile(file, header, "webui/public/images/conversation", userId)
	if err != nil {
		HandleError(w, NewAPIError("Error saving file", http.StatusBadRequest))
		return
	}

	payload := database.SendMessageRequest{
		ToUserId:  uint(toUserId),
		IsGroup:   isGroup,
		GroupId:   uint(groupId),
		PhotoPath: filepath,
	}

	result, err := h.SendMessage(userId, &payload)
	if err != nil {
		HandleError(w, NewAPIError("Failed to send message", http.StatusBadRequest))
		return
	}

	messageResponse, err := SendMessageResponseHandler(&payload, userId, result, false)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	messageResponse["photoPath"] = filepath

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messageResponse)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
