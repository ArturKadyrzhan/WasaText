package database

func (r *Repository) GetUser(user *User) (*User, error) {
	result := r.database.Model(&User{}).Where("username = ?", user.Username).Find(&user)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}
	return user, nil

}
