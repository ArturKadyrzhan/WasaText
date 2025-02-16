package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetGroupMessages(groupID uint) (*[]Message, error) {
	var conversation Conversation
	// Retrieve conversation by group ID
	query := `SELECT * FROM conversations WHERE group_id = ?`
	err := db.c.QueryRow(query, groupID).Scan(&conversation.ID, &conversation.User1ID, &conversation.User2ID, &conversation.GroupID, &conversation.IsGroup, &conversation.CreatedAt, &conversation.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &[]Message{}, nil
		}
		return &[]Message{}, fmt.Errorf("failed to fetch conversation: %w", err)
	}

	// Retrieve messages for this conversation
	var messages []Message
	messageQuery := `SELECT * FROM messages WHERE conversation_id = ? ORDER BY created_at ASC`
	rows, err := db.c.Query(messageQuery, conversation.ID)
	if err != nil {
		return &[]Message{}, fmt.Errorf("failed to fetch messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.ConversationID, &message.SenderID, &message.Content, &message.MessageType, &message.CreatedAt, &message.IsRead, &message.RepliedMessageID); err != nil {
			return &[]Message{}, fmt.Errorf("failed to scan message: %w", err)
		}

		var reaction Reaction
		err := db.c.QueryRow(`SELECT id, message_id, user_id, reaction, created_at FROM reactions WHERE message_id = ? LIMIT 1`, message.ID).
			Scan(&reaction.ID, &reaction.MessageID, &reaction.UserID, &reaction.Reaction, &reaction.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				message.Reactions = Reaction{}
			} else {
				return &[]Message{}, fmt.Errorf("failed to query reaction: %w", err)
			}
		} else {
			message.Reactions = reaction
		}
		user, err := db.GetUserById(&User{ID: message.SenderID})
		if err != nil {
			return &[]Message{}, fmt.Errorf("failed to get senders username: %w", err)
		}
		message.Sender = *user

		messages = append(messages, message)
	}
	return &messages, nil

}
