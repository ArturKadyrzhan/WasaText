package api

import (
	"WasaText/service/consts"
	"WasaText/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var convUserId database.GetMessagesRequest
	if err := json.NewDecoder(r.Body).Decode(&convUserId); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	var messages *[]database.Message
	var err error

	if !convUserId.IsGroup {
		messages, err = h.db.GetPrivateMessages(userId, convUserId.UserOrGroupId)
		if err != nil {
			HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
			return
		}
	} else {
		messages, err = h.db.GetGroupMessages(convUserId.UserOrGroupId)
		if err != nil {
			HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
			return
		}
	}

	var response []MessagesResponse

	if messages != nil {
		messageMap := make(map[uint]database.Message)
		for _, message := range *messages {
			messageMap[message.ID] = message
		}

		for _, message := range *messages {

			var emoji string
			repliedMessage := RepliedMessageResponse{}

			if message.Reactions.ID != 0 {
				emoji = message.Reactions.Reaction
			}

			if message.RepliedMessageID != nil {
				repliedMsg, found := messageMap[*message.RepliedMessageID]
				if found {
					repliedMessage = RepliedMessageResponse{
						MessageId: repliedMsg.ID,
						Content:   repliedMsg.Content,
					}
				}
			}

			if message.MessageType == consts.MessageTypeText {
				response = append(response, MessagesResponse{
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
				response = append(response, MessagesResponse{
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

	res := map[string][]MessagesResponse{"messages": response}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
