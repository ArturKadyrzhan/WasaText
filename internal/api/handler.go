package api

import (
	"WasaText/cmd/database/models"
	"WasaText/internal/helpers"
	"WasaText/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) doLogin(w http.ResponseWriter, r *http.Request) {
	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	//fmt.Println(input, "error zdes' tt")
	user, err := h.Service.CreateOrGetUser(&input)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Failed to retrieve token", http.StatusBadRequest))
		return
	}

	token, err := helpers.GenerateSessionToken(user)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Failed to retrieve token", http.StatusBadRequest))
		return
	}

	res := map[string]interface{}{
		"token":    token,
		"id":       user.ID,
		"username": user.Username,
	}

	fmt.Println(token)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Failed to encode response", http.StatusBadRequest))
		return
	}
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)
	query := r.URL.Query().Get("search")
	if query == "" {
		helpers.HandleError(w, helpers.NewAPIError("Query is empty", http.StatusBadRequest))
		return
	}

	users, err := h.Service.GetUsers(query, userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	//slice
	res := map[string][]models.User{"users": *users}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
	}

}

func (h *Handler) getConversations(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)
	info, err := h.Service.GetConversations(userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	res := map[string]interface{}{
		"result": info,
	}
	//
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

}

func (h *Handler) sendMessage(w http.ResponseWriter, r *http.Request) {
	var input helpers.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	userId := r.Context().Value("userId").(uint)

	result, err := h.Service.SendMessage(userId, &input)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	message, err := helpers.SendMessageResponseHandler(&input, userId, result, false)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}

func (h *Handler) getMessages(w http.ResponseWriter, r *http.Request) {
	var convUserId helpers.GetMessagesRequest
	if err := json.NewDecoder(r.Body).Decode(&convUserId); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	userId := r.Context().Value("userId").(uint)

	messages, err := h.Service.GetMessages(userId, &convUserId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	res := map[string][]helpers.MessagesResponse{"messages": *messages}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}

func (h *Handler) markAsRead(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)
	var input helpers.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	result, err := h.Service.MarkAsRead(userId, &input)
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

func (h *Handler) createGroup(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Unable to parse form data", http.StatusUnprocessableEntity))
		return
	}

	groupName := r.FormValue("groupName")
	if groupName == "" {
		helpers.HandleError(w, helpers.NewAPIError("Group name is required", http.StatusBadRequest))
		return
	}

	selectedUsers := r.FormValue("selectedUsers")
	if selectedUsers == "" {
		helpers.HandleError(w, helpers.NewAPIError("Selected users are required", http.StatusBadRequest))
		return
	}

	var users []models.User
	err = json.Unmarshal([]byte(selectedUsers), &users)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Invalid selected users data", http.StatusBadRequest))
		return
	}

	file, header, err := r.FormFile("profilePhoto")
	if err != nil && err != http.ErrMissingFile {
		helpers.HandleError(w, helpers.NewAPIError("Error retrieving file", http.StatusBadRequest))
		return
	}
	filepath, err := helpers.SaveUploadedFile(file, header, "webui/public/images/group", userId)

	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	result, err := h.Service.CreateGroups(&helpers.CreateGroupRequest{
		GroupName:      groupName,
		GroupPhotoPath: filepath,
		Users:          users,
	}, userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	res := map[string]bool{"success": result}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}

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
	if err != nil && err != http.ErrMissingFile {
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

	result, err := h.Service.SendMessage(userId, &payload)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Failed to send message", http.StatusBadRequest))
		return
	}

	message, err := helpers.SendMessageResponseHandler(&payload, userId, result, false)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	message["photoPath"] = filepath

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}

func parseOptionalInt(value string) (int, error) {
	if value == "undefined" || value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

func (h *Handler) uploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	file, header, err := r.FormFile("profile_picture")
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	defer file.Close()

	filepath, err := helpers.SaveUploadedFile(file, header, "webui/public/images/profile", userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	result, err := h.Service.UpdateUserProfile(userId, filepath)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	res := map[string]bool{"success": result}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}

func (h *Handler) addUsersToGroup(w http.ResponseWriter, r *http.Request) {
	var input helpers.AddUsersToGroup
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	userId := r.Context().Value("userId").(uint)

	result, err := h.Service.AddUsersToGroup(userId, &input)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	res := map[string]bool{"success": result}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}

func (h *Handler) deleteMessage(w http.ResponseWriter, r *http.Request) {
	var input helpers.DeleteMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	result, err := h.Service.DeleteMessage(input.MessageId)
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

func (h *Handler) leaveGroup(w http.ResponseWriter, r *http.Request) {
	var input helpers.GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	userId := r.Context().Value("userId").(uint)

	result, err := h.Service.LeaveGroup(userId, &input)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	res := map[string]bool{"success": result}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}

func (h *Handler) commentMessage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	var input helpers.CommentMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	result, err := h.Service.CommentMessage(&input, userId)
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

func (h *Handler) uncommentMessage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	var input helpers.UncommentMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	result, err := h.Service.UncommentMessage(&input, userId)
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

func (h *Handler) forwardMessage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	var input helpers.ForwardMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	result, err := h.Service.ForwardMessage(&input, userId)
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
