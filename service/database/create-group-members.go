package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) CreateGroupMembers(userId, addedById, groupId uint) (bool, error) {
	// Insert into the group_members table
	query := `
		INSERT INTO group_members (group_id, user_id, added_by, added_at)
		VALUES (?, ?, ?, ?)`
	_, err := db.c.Exec(query, groupId, userId, addedById, time.Now())
	if err != nil {
		return false, fmt.Errorf("failed to add group member: %w", err)
	}
	return true, nil
}
