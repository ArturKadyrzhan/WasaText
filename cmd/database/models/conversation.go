package models

import "time"

type Conversation struct {
	ID        uint  ` gorm:"primaryKey"`
	User1ID   *uint // NULL для групп
	User2ID   *uint // NULL для групп
	GroupID   *uint // NULL для личных бесед
	IsGroup   bool
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	Group Group `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE"`
	User1 User  `gorm:"foreignKey:User1ID;constraint:OnDelete:CASCADE"`
	User2 User  ` gorm:"foreignKey:User2ID;constraint:OnDelete:CASCADE"`
}

func (c *Conversation) TableName() string {
	return "conversations"
}
