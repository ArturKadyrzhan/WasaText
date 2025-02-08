package database

import "fmt"

func (r *Repository) CreateUser(user *User) (*User, error) {
	fmt.Println(user, "error etc")
	result := r.database.Create(user)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}
	return user, nil
}
