package database

import "fmt"

func (db *appdbimpl) UpdateGroupPhoto(id int, groupPhoto string) (bool, error) {

	query := `
		UPDATE groups 
		SET  group_photo_url = ?
		WHERE id = ?`

	result, err := db.c.Exec(query, groupPhoto, id)
	if err != nil {
		return false, fmt.Errorf("failed to update group photo: %w", err)
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
