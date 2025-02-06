package api

import (
	"WasaText/service/consts"
	"WasaText/service/database/models"
	"WasaText/service/helpers"
	"encoding/json"
	"net/http"
)

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
