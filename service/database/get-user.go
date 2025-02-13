package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) GetUser(user *User) (*User, error) {
	query := `SELECT * FROM users WHERE username = ?`
	err := db.c.QueryRow(query, user.Username).Scan(&user.ID, &user.Username, &user.ProfilePhotoURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &User{}, nil
		}
		return &User{}, fmt.Errorf("failed to fetch user: %w", err)
	}
	return user, nil
}
