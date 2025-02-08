package database

import (
	"gorm.io/gorm"
	"log"
)

func (r *Repository) GetGroupMessages(groupID uint) (*[]Message, error) {
	var conversation Conversation
	err := r.database.Where("group_id = ?", groupID).First(&conversation).Error
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
