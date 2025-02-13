package database

import "fmt"

func (db *appdbimpl) CreateGroup(payload *CreateGroupRequest, userId uint) (*Group, error) {
	var group Group
	group.CreatedBy = userId
	group.Name = payload.GroupName
	group.GroupPhotoURL = payload.GroupPhotoPath

	// Insert into the groups table
	query := `
		INSERT INTO groups (name, group_photo_url, created_by)
		VALUES (?, ?, ?)`
	result, err := db.c.Exec(query, group.Name, group.GroupPhotoURL, group.CreatedBy)
	if err != nil {
		return &Group{}, fmt.Errorf("failed to create group: %w", err)
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return &Group{}, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	group.ID = uint(lastInsertID)
	return &group, nil
}
