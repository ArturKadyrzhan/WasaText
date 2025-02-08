package database

import (
	"gorm.io/gorm"
	"log"
)

func (r *Repository) GetPrivateMessages(user1ID uint, user2ID uint) (*[]Message, error) {
	var conversation Conversation
	err := r.database.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", user1ID, user2ID, user2ID, user1ID).First(&conversation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println("Error fetching conversation:", err)
		return nil, err
	}

	var messages []Message
	err = r.database.
		Where("conversation_id = ?", conversation.ID).
		Order("created_at asc").
		Preload("Sender").
		Preload("Reactions").
		Find(&messages).Error

	return &messages, nil
}
