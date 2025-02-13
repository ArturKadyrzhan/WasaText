package database

import "fmt"

func (db *appdbimpl) CreateUser(user *User) (*User, error) {
	query := `
		INSERT INTO users (username, profile_photo_url)
		VALUES (?, ?)`
	result, err := db.c.Exec(query, user.Username, user.ProfilePhotoURL)
	if err != nil {
		return &User{}, fmt.Errorf("failed to create user: %w", err)
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return &User{}, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	user.ID = uint(lastInsertID)
	return user, nil
}
