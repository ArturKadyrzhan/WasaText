package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type AppDatabase interface {
	CheckGroupConversation(groupId uint) (*Conversation, error)
	CheckPrivateConversation(user1Id uint, user2Id uint) (*Conversation, error)
	CommentMessage(payload *CommentMessage, userId uint) (bool, error)
	CreateGroup(payload *CreateGroupRequest, userId uint) (*Group, error)
	CreateGroupMembers(userId, addedById, groupId uint) (bool, error)
	CreateMessage(message *Message) (*Message, error)
	CreateUser(user *User) (*User, error)
	DeleteGroupMember(userId uint, group Group) (bool, error)
	DeleteMessage(msgID uint) (bool, error)
	GetConversationsUsers(userId uint) (*[]User, error)
	GetGroupMembers(groupId uint) (*[]GroupMember, error)
	GetGroupMessages(groupID uint) (*[]Message, error)
	GetGroups(userId uint) (*[]Group, error)
	GetPrivateMessages(user1ID uint, user2ID uint) (*[]Message, error)
	GetUser(user *User) (*User, error)
	GetUsers(query string, userId uint) (*[]User, error)
	MarkAsRead(convId, userId uint) (bool, error)
	UncommentMessage(payload *UncommentMessage, userId uint) (bool, error)
	UpdateUserProfile(user *User) (bool, error)
	GetUserById(user *User) (*User, error)
	UpdateGroupPhoto(id int, groupPhoto string) (bool, error)
	UpdateGroupName(id int, groupName string) (bool, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func createDatabase(db *sql.DB) error {
	tables := [7]string{

		`CREATE TABLE IF NOT EXISTS users (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(16) NOT NULL UNIQUE,
			profile_photo_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME
		);`,

		`CREATE TABLE IF NOT EXISTS groups (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			group_photo_url TEXT,
			created_by INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME,
			FOREIGN KEY (created_by) REFERENCES users(ID) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS conversations (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			user1_id INTEGER,
			user2_id INTEGER,
			group_id INTEGER,
			is_group BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME,
			FOREIGN KEY (user1_id) REFERENCES users(ID) ON DELETE CASCADE,
			FOREIGN KEY (user2_id) REFERENCES users(ID) ON DELETE CASCADE,
			FOREIGN KEY (group_id) REFERENCES groups(ID) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS group_members (
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			added_by INTEGER NOT NULL,
			added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (group_id, user_id),
			FOREIGN KEY (group_id) REFERENCES groups(ID) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE,
			FOREIGN KEY (added_by) REFERENCES users(ID) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS messages (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			conversation_id INTEGER NOT NULL,
			sender_id INTEGER NOT NULL,
			forwarded_by_id INTEGER,
			content TEXT NOT NULL,
			message_type VARCHAR(10) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_read BOOLEAN DEFAULT FALSE,
			replied_message_id INTEGER,
			FOREIGN KEY (conversation_id) REFERENCES conversations(ID) ON DELETE CASCADE,
			FOREIGN KEY (sender_id) REFERENCES users(ID) ON DELETE CASCADE,
			FOREIGN KEY (forwarded_by_id) REFERENCES users(ID) ON DELETE CASCADE,
			FOREIGN KEY (replied_message_id) REFERENCES messages(ID) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS reactions (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			message_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			reaction TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (message_id) REFERENCES messages(ID) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
		);`,
	}

	for _, sqlStmt := range tables {
		_, err := db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("error creating table: %w", err)
		}
	}

	return nil
}
