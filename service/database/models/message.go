package models

import "time"

type Message struct {
	ID               uint      `gorm:"primaryKey"`
	ConversationID   uint      `gorm:"not null"`
	SenderID         uint      `gorm:"not null"`
	Content          string    `gorm:"type:text"`
	MessageType      string    `gorm:"not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	IsRead           bool      `gorm:"default:false"`
	RepliedMessageID *uint

	Conversation Conversation `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE"`
	Sender       User         `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	Reactions    Reaction     `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE"`
}
