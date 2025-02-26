package database

import "fmt"

func (db *appdbimpl) GetGroups(userId uint) (*[]Group, error) {
	var groups []Group

	query := `
		SELECT DISTINCT g.* 
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = ?`

	rows, err := db.c.Query(query, userId)
	if err != nil {
		return &[]Group{}, fmt.Errorf("failed to get groups: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.ID, &group.Name, &group.GroupPhotoURL, &group.CreatedBy, &group.CreatedAt, &group.UpdatedAt); err != nil {
			return &[]Group{}, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return &[]Group{}, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return &groups, nil
}
