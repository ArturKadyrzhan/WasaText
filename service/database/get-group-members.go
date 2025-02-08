package database

func (r *Repository) GetGroupMembers(groupId uint) (*[]GroupMember, error) {
	var groupMember []GroupMember
	result := r.database.Model(&GroupMember{}).
		Where("group_id = ?", groupId).Preload("User").
		Find(&groupMember)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}

	return &groupMember, nil
}
