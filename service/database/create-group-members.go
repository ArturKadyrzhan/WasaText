package database

import "time"

func (r *Repository) CreateGroupMembers(userId, addedById, groupId uint) (bool, error) {
	var groupMember GroupMember
	groupMember.GroupID = groupId
	groupMember.UserID = userId
	groupMember.AddedBy = addedById
	groupMember.AddedAt = time.Now()

	result := r.database.Create(&groupMember)
	if result.Error != nil {
		msg := result.Error
		return false, msg
	}

	return true, nil
}
