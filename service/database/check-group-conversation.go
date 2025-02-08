package database

import (
	"errors"
	"gorm.io/gorm"
)

func (r *Repository) CheckGroupConversation(groupId uint) (*Conversation, error) {
	var conv Conversation

	result := r.database.Model(&Conversation{}).
		Where("group_id = ?", groupId).
		Where("is_group = ?", true).
		First(&conv)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			conv = Conversation{
				GroupID: &groupId,
				IsGroup: true,
			}

			return r.CreateConversation(&conv)
		}
		return nil, result.Error
	}

	return &conv, nil
}
