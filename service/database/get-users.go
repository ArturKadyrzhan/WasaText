package database

import "strings"

func (r *Repository) GetUsers(query string, userId uint) (*[]User, error) {
	var users []User
	if err := r.database.Where("LOWER(username) LIKE ? AND id != ?",
		"%"+strings.ToLower(query)+"%", userId).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}
