package database

import (
	"fmt"
	"strings"
)

func (db *appdbimpl) GetUsers(query string, userId uint) (*[]User, error) {
	users := []User{}
	searchQuery := `SELECT * FROM users WHERE LOWER(username) LIKE ? AND id != ?`
	rows, err := db.c.Query(searchQuery, "%"+strings.ToLower(query)+"%", userId)
	if err != nil {
		return &users, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.ProfilePhotoURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return &users, fmt.Errorf("failed to scan user: %w", err)
		}
		fmt.Println(user, "get-users database-dan")
		users = append(users, user)
	}
	return &users, nil
}
