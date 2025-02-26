package database

import "fmt"

func (db *appdbimpl) UpdateUserProfile(user *User) (bool, error) {

	query := `
		UPDATE users 
		SET profile_photo_url = ? , username =?
		WHERE id = ?`

	result, err := db.c.Exec(query, user.ProfilePhotoURL, user.Username, user.ID)
	if err != nil {
		return false, fmt.Errorf("failed to update user profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
