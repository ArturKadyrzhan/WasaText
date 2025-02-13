package database

import "fmt"

func (db *appdbimpl) GetConversationsUsers(userId uint) (*[]User, error) {
	var users []User
	query := `
		SELECT DISTINCT u.*
		FROM users u
		JOIN conversations c ON (u.id = c.user1_id OR u.id = c.user2_id)
		WHERE (c.user1_id = ? OR c.user2_id = ?)
		AND u.id != ?`
	rows, err := db.c.Query(query, userId, userId, userId)
	if err != nil {
		return &[]User{}, fmt.Errorf("failed to fetch conversation users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.ProfilePhotoURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return &[]User{}, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return &users, nil
}
