package database

import (
	"fmt"
)

func (db *appdbimpl) GetGroupMembers(groupId uint) (*[]GroupMember, error) {
	var groupMembers []GroupMember
	query := `
		SELECT gm.group_id, gm.user_id, gm.added_by, gm.added_at, g.name, u.username
		FROM group_members gm
		JOIN groups g ON gm.group_id = g.id
		JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = ?`
	rows, err := db.c.Query(query, groupId)
	if err != nil {
		return &[]GroupMember{}, fmt.Errorf("failed to fetch group members: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var groupMember GroupMember
		if err := rows.Scan(&groupMember.GroupID, &groupMember.UserID, &groupMember.AddedBy, &groupMember.AddedAt, &groupMember.Group.Name, &groupMember.User.Username); err != nil {
			return &[]GroupMember{}, fmt.Errorf("failed to scan group member: %w", err)
		}
		groupMembers = append(groupMembers, groupMember)
	}

	if err := rows.Err(); err != nil {
		return &[]GroupMember{}, fmt.Errorf("error iterating over group members: %w", err)
	}
	return &groupMembers, nil
}
