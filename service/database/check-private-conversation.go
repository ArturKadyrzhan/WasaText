package database

import (
	"errors"
	"gorm.io/gorm"
)

func (r *Repository) CheckPrivateConversation(user1Id uint, user2Id uint) (*Conversation, error) {
	var conv Conversation

	result := r.database.Model(Conversation{}).
		Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", user1Id, user2Id, user2Id, user1Id).
		Where("is_group = ?", false).
		First(&conv)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			user1 := User{ID: user1Id}
			user2 := User{ID: user2Id}

			conv = Conversation{
				User1ID: &user1.ID,
				User2ID: &user2.ID,
			}

			return r.CreateConversation(&conv)
		}
		return nil, result.Error
	}

	return &conv, nil
}
