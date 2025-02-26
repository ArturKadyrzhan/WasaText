package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) CheckGroupConversation(groupId uint) (*Conversation, error) {
	var conv Conversation
	query := `SELECT id, user1_id, user2_id, group_id, is_group,created_at,updated_at FROM conversations WHERE group_id = ? AND is_group = ? LIMIT 1`
	err := db.c.QueryRow(query, groupId, true).Scan(&conv.ID, &conv.User1ID, &conv.User2ID, &conv.GroupID, &conv.IsGroup,
		&conv.CreatedAt, &conv.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			conv = Conversation{

				GroupID: &groupId,
				IsGroup: true,
			}
			return db.CreateConversation(&conv)
		}
		return &Conversation{}, fmt.Errorf("failed to check group conversation: %w", err)
	}
	return &conv, nil
}
