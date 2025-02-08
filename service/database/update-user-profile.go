package database

func (r *Repository) UpdateUserProfile(user *User) (bool, error) {
	result := r.database.Model(&user).Where("id = ?", user.ID).Update("profile_photo_url", user.ProfilePhotoURL)
	if result.Error != nil {
		msg := result.Error
		return false, msg
	}

	return true, nil
}
