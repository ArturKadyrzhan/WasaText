package database

import (
	"fmt"
)

func (db *appdbimpl) DeleteMessage(msgID uint) (bool, error) {
	query := `
		DELETE FROM messages
		WHERE id = ?`
	result, err := db.c.Exec(query, msgID)
	if err != nil {
		return false, fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}
