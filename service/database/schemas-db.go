package database

import "time"

type Conversation struct {
	ID        uint  ` gorm:"primaryKey"`
	User1ID   *uint // NULL for groups
	User2ID   *uint // NULL for local conversations
	GroupID   *uint // NULL for local conversations
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

type Reaction struct {
	ID        uint      ` gorm:"primaryKey"`
	MessageID uint      `  gorm:"not null"`
	UserID    uint      ` gorm:"not null"`
	Reaction  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (r *Reaction) TableName() string {
	return "reactions"
}

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
