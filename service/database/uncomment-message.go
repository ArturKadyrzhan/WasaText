package database

import "fmt"

func (db *appdbimpl) UncommentMessage(payload *UncommentMessage, userId uint) (bool, error) {
	// Delete the reaction from the reactions table to remove the comment (emoji)
	query := `
		DELETE FROM reactions
		WHERE message_id = ? AND user_id = ?`
	result, err := db.c.Exec(query, payload.MessageId, userId)
	if err != nil {
		return false, fmt.Errorf("failed to uncomment message: %w", err)
	}

	// Check if the reaction was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to check rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}
