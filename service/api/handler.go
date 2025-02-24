package api

import (
	"WasaText/service/consts"
	"WasaText/service/database"
	"strconv"
	"time"
)

func parseOptionalInt(value string) (int, error) {
	if value == "undefined" || value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

func (h *_router) SendMessage(userId uint, input *database.SendMessageRequest) (*database.Message, error) {
	var conv *database.Conversation
	var err error

	if !input.IsGroup {
		conv, err = h.db.CheckPrivateConversation(userId, input.ToUserId)
		if err != nil {
			return nil, err
		}
	} else {
		conv, err = h.db.CheckGroupConversation(input.GroupId)
		if err != nil {
			return nil, err
		}
	}
	var message database.Message

	if input.PhotoPath == "" {
		message = database.Message{
			ConversationID:   conv.ID,
			SenderID:         userId,
			Content:          input.Text,
			MessageType:      consts.MessageTypeText,
			CreatedAt:        time.Now(),
			RepliedMessageID: input.RepliedMessageId,
			ForwardedBy:      database.User{ID: uint(input.ForwardedById)},
		}
	} else {
		message = database.Message{
			ConversationID:   conv.ID,
			SenderID:         userId,
			Content:          input.PhotoPath,
			MessageType:      consts.MessageTypePhoto,
			CreatedAt:        time.Now(),
			RepliedMessageID: input.RepliedMessageId,
			ForwardedBy:      database.User{ID: uint(input.ForwardedById)},
		}
	}

	return h.db.CreateMessage(&message)
}
