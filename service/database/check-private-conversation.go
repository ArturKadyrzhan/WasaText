package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) CheckPrivateConversation(user1Id uint, user2Id uint) (*Conversation, error) {
	var conv Conversation

	// Query to check if a private conversation exists between the two users
	query := `SELECT id, user1_id, user2_id, is_group 
			  FROM conversations 
			  WHERE ((user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)) 
			  AND is_group = ? 
			  LIMIT 1`

	err := db.c.QueryRow(query, user1Id, user2Id, user2Id, user1Id, false).
		Scan(&conv.ID, &conv.User1ID, &conv.User2ID, &conv.IsGroup)

	if err != nil {
		// If no record is found, create a new private conversation
		if errors.Is(err, sql.ErrNoRows) {
			conv = Conversation{
				User1ID: &user1Id,
				User2ID: &user2Id,
				IsGroup: false,
			}

			// Call the CreateConversation method to insert the new conversation
			return db.CreateConversation(&conv)
		}
		// Return other possible errors from the query
		return &Conversation{}, fmt.Errorf("failed to check private conversation: %w", err)
	}

	// Return the existing conversation
	return &conv, nil
}
