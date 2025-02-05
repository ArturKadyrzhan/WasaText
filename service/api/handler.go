package api

import (
	"WasaText/service/consts"
	"WasaText/service/database/models"
	"WasaText/service/helpers"
	"WasaText/service/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"
)

type Handler struct {
	Repository *repositories.Repository
}

func NewHandler(repository *repositories.Repository) *Handler {
	return &Handler{Repository: repository}
}

func (h *Handler) doLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	//fmt.Println(input, "error zdes' tt")
	existUser, err := h.Repository.GetUser(&user)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Failed to retrieve token", http.StatusBadRequest))
		return
	}
	var createdUser *models.User
	if existUser.ID == 0 {
		createdUser, err = h.Repository.CreateUser(&user)
		if err != nil {
			helpers.HandleError(w, helpers.NewAPIError("Failed to retrieve token", http.StatusBadRequest))
			return
		}

	}

	token, err := helpers.GenerateSessionToken(createdUser)
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

	users, err := h.Repository.GetUsers(query, userId)
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
		return
	}

}

func (h *Handler) getConversations(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)
	info := make(map[string]interface{})
	users, err := h.Repository.GetConversationsUsers(userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	groups, err := h.Repository.GetGroups(userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	fmt.Println(groups, "!!!! ERROR")
	info["users"] = users
	info["groups"] = groups

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

	result, err := h.SendMessage(userId, &input)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	messageResponse, err := helpers.SendMessageResponseHandler(&input, userId, result, false)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messageResponse)
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

	var messages *[]models.Message
	var err error

	if !convUserId.IsGroup {
		messages, err = h.Repository.GetPrivateMessages(userId, convUserId.UserOrGroupId)
		if err != nil {
			helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
			return
		}
	} else {
		messages, err = h.Repository.GetGroupMessages(convUserId.UserOrGroupId)
		if err != nil {
			helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
			return
		}
	}

	var response []helpers.MessagesResponse

	if messages != nil {
		messageMap := make(map[uint]models.Message)
		for _, message := range *messages {
			messageMap[message.ID] = message
		}

		for _, message := range *messages {

			var emoji string
			repliedMessage := helpers.RepliedMessageResponse{}

			if message.Reactions.ID != 0 {
				emoji = message.Reactions.Reaction
			}

			if message.RepliedMessageID != nil {
				repliedMsg, found := messageMap[*message.RepliedMessageID]
				if found {
					repliedMessage = helpers.RepliedMessageResponse{
						MessageId: repliedMsg.ID,
						Content:   repliedMsg.Content,
					}
				}
			}

			if message.MessageType == consts.MessageTypeText {
				response = append(response, helpers.MessagesResponse{
					MessageId:      message.ID,
					Message:        message.Content,         // Message text
					UserId:         message.SenderID,        // Sender user ID
					ConversationId: message.ConversationID,  // Conversation ID
					Username:       message.Sender.Username, // Conversation ID
					CreatedAt:      message.CreatedAt,       // Conversation ID
					IsPhoto:        false,
					IsRead:         message.IsRead,
					Emoji:          emoji,
					RepliedMessage: repliedMessage,
				})
			} else {
				response = append(response, helpers.MessagesResponse{
					MessageId:      message.ID,
					Message:        message.Content,         // Message text
					UserId:         message.SenderID,        // Sender user ID
					ConversationId: message.ConversationID,  // Conversation ID
					Username:       message.Sender.Username, // Conversation ID
					CreatedAt:      message.CreatedAt,       // Conversation ID
					IsPhoto:        true,
					IsRead:         message.IsRead,
					Emoji:          emoji,
					RepliedMessage: repliedMessage,
				})
			}
		}
	}

	res := map[string][]helpers.MessagesResponse{"messages": response}

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

	var conv *models.Conversation
	var err error

	if !input.IsGroup {
		conv, err = h.Repository.CheckPrivateConversation(userId, input.ToUserId)
		if err != nil {
			helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
			return
		}
	} else {
		conv, err = h.Repository.CheckGroupConversation(input.GroupId)
		if err != nil {
			helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
			return
		}
	}

	result, err := h.Repository.MarkAsRead(conv.ID, userId)
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
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError("Error retrieving file", http.StatusBadRequest))
		return
	}
	filepath, err := helpers.SaveUploadedFile(file, header, "webui/public/images/group", userId)

	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
	payload := &helpers.CreateGroupRequest{
		GroupName:      groupName,
		GroupPhotoPath: filepath,
		Users:          users,
	}
	group, err := h.Repository.CreateGroup(payload, userId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	_, err = h.Repository.CreateGroupMembers(userId, userId, group.ID)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	for _, user := range payload.Users {
		_, err = h.Repository.CreateGroupMembers(user.ID, userId, group.ID)
		if err != nil {
			helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
			return
		}
	}

	res := map[string]bool{"success": true}

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

	var user models.User
	user.ID = userId
	user.ProfilePhotoURL = filepath
	result, err := h.Repository.UpdateUserProfile(&user)
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

	users, err := h.Repository.GetGroupMembers(input.GroupId)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	var userIds = make([]uint, len(*users))

	for _, user := range *users {
		userIds = append(userIds, user.UserID)
	}

	for _, v := range input.Users {
		if !slices.Contains(userIds, v.ID) {
			_, err = h.Repository.CreateGroupMembers(v.ID, userId, input.GroupId)
			if err != nil {
				helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
				return
			}
		}
	}

	res := map[string]bool{"success": true}

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

func (h *Handler) leaveGroup(w http.ResponseWriter, r *http.Request) {
	var input helpers.GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	userId := r.Context().Value("userId").(uint)

	result, err := h.Repository.DeleteGroupMember(userId, input.Group)
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

func (h *Handler) uncommentMessage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	var input helpers.UncommentMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	result, err := h.Repository.UncommentMessage(&input, userId)
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

func (h *Handler) SendMessage(userId uint, input *helpers.SendMessageRequest) (*models.Message, error) {
	var conv *models.Conversation
	var err error

	if !input.IsGroup {
		conv, err = h.Repository.CheckPrivateConversation(userId, input.ToUserId)
		if err != nil {
			return nil, err
		}
	} else {
		conv, err = h.Repository.CheckGroupConversation(input.GroupId)
		if err != nil {
			return nil, err
		}
	}
	var message models.Message

	if input.PhotoPath == "" {
		message = models.Message{
			ConversationID:   conv.ID,
			SenderID:         userId,
			Content:          input.Text,
			MessageType:      consts.MessageTypeText,
			CreatedAt:        time.Now(),
			RepliedMessageID: input.RepliedMessageId,
		}
	} else {
		message = models.Message{
			ConversationID:   conv.ID,
			SenderID:         userId,
			Content:          input.PhotoPath,
			MessageType:      consts.MessageTypePhoto,
			CreatedAt:        time.Now(),
			RepliedMessageID: input.RepliedMessageId,
		}
	}

	return h.Repository.CreateMessage(&message)
}

func (h *Handler) forwardMessage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	var input helpers.ForwardMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
		return
	}
	if input.IsPhoto {
		for _, user := range input.Users {
			sendMessagePayload := helpers.SendMessageRequest{
				ToUserId:  user.ID,
				IsGroup:   false,
				PhotoPath: input.Text,
			}
			_, err := h.SendMessage(userId, &sendMessagePayload)
			if err != nil {
				helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
				return
			}
		}
	} else {
		for _, user := range input.Users {
			sendMessagePayload := helpers.SendMessageRequest{
				ToUserId: user.ID,
				IsGroup:  false,
				Text:     input.Text,
			}

			_, err := h.SendMessage(userId, &sendMessagePayload)
			if err != nil {
				helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusUnprocessableEntity))
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(true)
	if err != nil {
		helpers.HandleError(w, helpers.NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
