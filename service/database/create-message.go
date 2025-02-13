package database

import "fmt"

func (db *appdbimpl) CreateMessage(message *Message) (*Message, error) {
	// Insert into the messages table
	query := `
		INSERT INTO messages (conversation_id, sender_id, content, message_type, is_read, replied_message_id)
		VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.c.Exec(query, message.ConversationID, message.SenderID, message.Content, message.MessageType, message.IsRead, message.RepliedMessageID)
	if err != nil {
		return &Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return &Message{}, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	message.ID = uint(lastInsertID)
	return message, nil
}

func (db *appdbimpl) CreateConversation(conv *Conversation) (*Conversation, error) {
	// Insert into the conversations table
	query := `
		INSERT INTO conversations (user1_id, user2_id, group_id, is_group)
		VALUES (?, ?, ?, ?)`
	result, err := db.c.Exec(query, conv.User1ID, conv.User2ID, conv.GroupID, conv.IsGroup)
	if err != nil {
		return &Conversation{}, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return &Conversation{}, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	conv.ID = uint(lastInsertID)
	return conv, nil
}
