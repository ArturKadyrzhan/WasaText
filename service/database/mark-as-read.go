package database

import "fmt"

func (db *appdbimpl) MarkAsRead(convId, userId uint) (bool, error) {
	query := `
		UPDATE messages
		SET is_read = ?
		WHERE conversation_id = ? AND is_read = ? AND sender_id != ?`
	result, err := db.c.Exec(query, true, convId, false, userId)
	if err != nil {
		return false, fmt.Errorf("failed to mark messages as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}
