package database

import (
	"log"
)

func (r *Repository) UncommentMessage(payload *UncommentMessage, userId uint) (bool, error) {
	res := r.database.Model(&Reaction{}).
		Where("message_id = ? AND user_id =?", payload.MessageId, userId).
		Delete(&Reaction{})
	if res.Error != nil {
		log.Println("Error updating is_read:", res.Error)
		return false, res.Error
	}

	if res.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
