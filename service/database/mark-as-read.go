package database

import "log"

func (r *Repository) MarkAsRead(convId, userId uint) (bool, error) {
	err := r.database.Model(&Message{}).
		Where("conversation_id = ? AND is_read = ? AND sender_id != ?", convId, false, userId).
		Update("is_read", true).Error
	if err != nil {
		log.Println("Error updating is_read:", err)
		return false, err
	}

	return true, nil
}
