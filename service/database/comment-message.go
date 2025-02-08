package database

import (
	"WasaText/service/api"
	"errors"
	"gorm.io/gorm"
)

func (r *Repository) CommentMessage(payload *api.CommentMessage, userId uint) (bool, error) {
	var reaction Reaction

	result := r.database.Where("message_id = ? AND user_id = ?", payload.MessageId, userId).First(&reaction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			reaction.MessageID = payload.MessageId
			reaction.UserID = userId
			reaction.Reaction = payload.Emoji

			if err := r.database.Create(&reaction).Error; err != nil {
				return false, err
			}
		}
		return false, result.Error

	}

	if err := r.database.Model(&Reaction{}).Where("message_id = ? AND user_id = ?", payload.MessageId, userId).Update("reaction", payload.Emoji).Error; err != nil {
		return false, err
	}

	return true, nil
}
