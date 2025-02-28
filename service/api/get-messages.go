package api

import (
	"WasaText/service/consts"   //constants for message types
	"WasaText/service/database" // database operations
	"encoding/json"
	"github.com/julienschmidt/httprouter" //  library for handling HTTP routes.
	"net/http"                            // for handling HTTP requests and responses
)

// h - function belongs to a router struct , w - send data back to client, r - incoming http request,
func (h *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Extracting Conversation Details from Request Body
	var convUserId database.GetMessagesRequest
	if err := json.NewDecoder(r.Body).Decode(&convUserId); err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}

	// Extracting the User ID from Request Context
	userVal := r.Context().Value(keyUserID)
	userId, ok := userVal.(uint)
	if !ok {
		HandleError(w, NewAPIError("invalid user ID in context", http.StatusInternalServerError))
		return
	}

	// Fetching Messages from the Database
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

	// Processing Messages for Response
	var response []MessagesResponse

	// Creates a messageMap to store messages by their IDs.
	// Loops through messages to: extract emojis,identify replied messages
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
			//// Building the Final Response, Adds each message to the response list.
			if message.MessageType == consts.MessageTypeText {
				response = append(response, MessagesResponse{
					MessageId:           message.ID,
					Message:             message.Content,
					UserId:              message.SenderID,
					ConversationId:      message.ConversationID,
					Username:            message.Sender.Username,
					CreatedAt:           message.CreatedAt,
					IsPhoto:             false,
					IsRead:              message.IsRead,
					Emoji:               emoji,
					RepliedMessage:      repliedMessage,
					ForwardedByUsername: message.ForwardedBy.Username,
				})
			} else {
				response = append(response, MessagesResponse{
					MessageId:           message.ID,
					Message:             message.Content,
					UserId:              message.SenderID,
					ConversationId:      message.ConversationID,
					Username:            message.Sender.Username,
					CreatedAt:           message.CreatedAt,
					IsPhoto:             true,
					IsRead:              message.IsRead,
					Emoji:               emoji,
					RepliedMessage:      repliedMessage,
					ForwardedByUsername: message.ForwardedBy.Username,
				})
			}
		}
	}

	// sending response
	res := map[string][]MessagesResponse{"messages": response}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, NewAPIError(err.Error(), http.StatusBadRequest))
		return
	}
}
