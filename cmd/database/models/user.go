package models

import "time"

type User struct {
	ID              uint       `gorm:"primaryKey"`
	Username        string     `gorm:"unique;not null"`
	ProfilePhotoURL string     `gorm:"type:text"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}
