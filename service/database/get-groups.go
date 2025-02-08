package database

func (r *Repository) GetGroups(userId uint) (*[]Group, error) {
	var groups []Group

	if err := r.database.Raw(`
        SELECT DISTINCT g.*
        FROM groups g
        JOIN group_members gm ON g.id = gm.group_id
        WHERE gm.user_id = ?`, userId).Scan(&groups).Error; err != nil {
		return nil, err
	}

	return &groups, nil
}
