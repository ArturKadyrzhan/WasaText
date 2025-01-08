package models

import "time"

type Group struct {
	ID            uint       `gorm:"primaryKey"`
	Name          string     ` gorm:"not null"`
	GroupPhotoURL string     ` gorm:"type:text"`
	CreatedBy     uint       `gorm:"not null"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`

	Creator User `gorm:"foreignKey:CreatedBy;constraint:OnDelete:CASCADE"`
}

func (g *Group) TableName() string {
	return "groups"
}
