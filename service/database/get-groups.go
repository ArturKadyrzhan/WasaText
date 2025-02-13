package database

import "fmt"

func (db *appdbimpl) GetGroups(userId uint) (*[]Group, error) {
	var groups []Group

	// Define the SELECT query to get the groups where the user is a member
	query := `
		SELECT DISTINCT g.* 
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = ?`

	// Execute the query
	rows, err := db.c.Query(query, userId)
	if err != nil {
		return &[]Group{}, fmt.Errorf("failed to get groups: %w", err)
	}
	defer rows.Close()

	// Iterate over the rows and scan them into the groups slice
	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.ID, &group.Name, &group.GroupPhotoURL, &group.CreatedBy, &group.CreatedAt, &group.UpdatedAt); err != nil {
			return &[]Group{}, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return &[]Group{}, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return &groups, nil
}
