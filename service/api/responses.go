package api

import (
	"WasaText/service/database"
	"time"
)

type MessagesResponse struct {
	MessageId      uint                   ` json:"message_id"`
	Message        string                 ` json:"message"`
	IsPhoto        bool                   `json:"is_photo"`
	UserId         uint                   `json:"user_id"`
	ConversationId uint                   `json:"conversation_id"`
	Username       string                 `json:"username"`
	IsRead         bool                   `json:"is_read"`
	Emoji          string                 `json:"emoji"`
	RepliedMessage RepliedMessageResponse `json:"replied_message"`
	CreatedAt      time.Time              `json:"createdAt"`
}

func SendMessageResponseHandler(input *database.SendMessageRequest, userId uint, result *database.Message, resultRead bool) (map[string]interface{}, error) {
	var fromUsername string
	message := map[string]interface{}{}

	if input.IsGroup {
		if input.PhotoPath != "" {
			message = map[string]interface{}{
				"id":         result.ID,
				"message":    input.PhotoPath,
				"isPhoto":    true,
				"sender":     userId,
				"groupId":    input.GroupId,
				"username":   fromUsername,
				"isReceived": true,
				"isRead":     resultRead,
				"createdAt":  time.Now(),
			}
		} else {
			message = map[string]interface{}{
				"id":         result.ID,
				"message":    input.Text,
				"isPhoto":    false,
				"sender":     userId,
				"groupId":    input.GroupId,
				"username":   fromUsername,
				"createdAt":  time.Now(),
				"isReceived": true,
				"isRead":     resultRead,
			}
		}
	} else {
		if input.PhotoPath != "" {
			message = map[string]interface{}{
				"id":         result.ID,
				"message":    input.PhotoPath,
				"isPhoto":    true,
				"sender":     userId,
				"createdAt":  time.Now(),
				"isReceived": true,
				"isRead":     resultRead,
			}
		} else {
			message = map[string]interface{}{
				"id":         result.ID,
				"message":    input.Text,
				"isPhoto":    false,
				"sender":     userId,
				"createdAt":  time.Now(),
				"isReceived": true,
				"isRead":     resultRead,
			}
		}
	}

	return message, nil
}

type RepliedMessageResponse struct {
	MessageId uint   `json:"message_id"`
	Content   string `json:"content"`
	Username  string `json:"username"`
}
