package database

import "fmt"

func (db *appdbimpl) CommentMessage(payload *CommentMessage, userId uint) (bool, error) {
	// Insert a new reaction into the reactions table to represent the comment (emoji)
	query := `
		INSERT INTO reactions (message_id, user_id, reaction)
		VALUES (?, ?, ?)`
	result, err := db.c.Exec(query, payload.MessageId, userId, payload.Emoji)
	if err != nil {
		return false, fmt.Errorf("failed to comment on message: %w", err)
	}

	// Check if the reaction was added
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to check rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}
