package database

func (r *Repository) DeleteGroupMember(userId uint, group Group) (bool, error) {
	result := r.database.Model(&GroupMember{}).
		Where("user_id = ? AND group_id = ?", userId, group.ID).
		Delete(&GroupMember{})

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
