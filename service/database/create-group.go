package database

func (r *Repository) CreateGroup(payload *CreateGroupRequest, userId uint) (*Group, error) {
	var group Group

	group.CreatedBy = userId
	group.Name = payload.GroupName
	group.GroupPhotoURL = payload.GroupPhotoPath

	result := r.database.Create(&group)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}

	return &group, nil
}
