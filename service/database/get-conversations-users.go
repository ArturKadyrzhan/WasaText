package database

func (r *Repository) GetConversationsUsers(userId uint) (*[]User, error) {
	var users []User

	if err := r.database.Raw(`
        SELECT DISTINCT u.*FROM users u  JOIN conversations c ON (u.id = c.user1_id OR u.id = c.user2_id)
        WHERE (c.user1_id = ? OR c.user2_id = ?)
          AND u.id != ?`, userId, userId, userId).Scan(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}
