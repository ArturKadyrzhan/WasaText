package database

import "fmt"

func (db *appdbimpl) UpdateUserProfile(user *User) (bool, error) {
	// Define the UPDATE query to update the profile_photo_url for the user
	query := `
		UPDATE users 
		SET profile_photo_url = ? , username =?
		WHERE id = ?`
	// Execute the UPDATE query
	result, err := db.c.Exec(query, user.ProfilePhotoURL, user.Username, user.ID)
	if err != nil {
		return false, fmt.Errorf("failed to update user profile: %w", err)
	}

	// Check how many rows were affected by the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to check affected rows: %w", err)
	}

	// If no rows were affected, return false (no user was found to update)
	if rowsAffected == 0 {
		return false, nil
	}

	// Return true indicating that the profile photo URL was successfully updated
	return true, nil
}
