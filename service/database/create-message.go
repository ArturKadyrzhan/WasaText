package database

func (r *Repository) CreateMessage(message *Message) (*Message, error) {
	result := r.database.Create(&message).Find(&message)
	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}

func (r *Repository) CreateConversation(conv *Conversation) (*Conversation, error) {
	result := r.database.Create(conv)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}

	return conv, nil
}
