package database

import "fmt"

func (db *appdbimpl) DeleteGroupMember(userId uint, group Group) (bool, error) {
	query := `
		DELETE FROM group_members
		WHERE user_id = ? AND group_id = ?`
	result, err := db.c.Exec(query, userId, group.ID)
	if err != nil {
		return false, fmt.Errorf("failed to delete group member: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}
