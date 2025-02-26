package database

import "fmt"

func (db *appdbimpl) UpdateGroupName(id int, groupName string) (bool, error) {

	query := `
		UPDATE groups 
		SET  name = ?
		WHERE id = ?`

	result, err := db.c.Exec(query, groupName, id)
	if err != nil {
		return false, fmt.Errorf("failed to update group name: %w", err)
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
