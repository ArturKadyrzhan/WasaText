package models

import "time"

type GroupMember struct {
	GroupID uint      `gorm:"primaryKey"`
	UserID  uint      `gorm:"primaryKey"`
	AddedBy uint      `gorm:"not null"`
	AddedAt time.Time `gorm:"autoCreateTime"`

	Group Group `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE"`
	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Adder User  `gorm:"foreignKey:AddedBy;constraint:OnDelete:CASCADE"`
}

func (gm *GroupMember) TableName() string {
	return "group_members"
}
