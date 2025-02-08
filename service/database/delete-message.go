package database

import "log"

func (r *Repository) DeleteMessage(msgID uint) (bool, error) {
	err := r.database.Model(&Message{}).
		Where("id = ?", msgID).
		Delete(&Message{}).Error
	if err != nil {
		log.Println("Error updating is_read:", err)
		return false, err
	}

	return true, nil
}
