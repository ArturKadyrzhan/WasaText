package api

import (
	"WasaText/service/consts"
	"WasaText/service/database"
	"WasaText/service/helpers"
	"strconv"
	"time"
)

type Handler struct {
	Repository *database.Repository
}

func parseOptionalInt(value string) (int, error) {
	if value == "undefined" || value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
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
