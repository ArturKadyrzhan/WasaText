/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// AppDatabase is the high level interface for the DB
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

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
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
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(16) NOT NULL UNIQUE,
			profile_photo_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME
		);`,

		// Groups table
		`CREATE TABLE IF NOT EXISTS groups (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			group_photo_url TEXT,
			created_by INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME,
			FOREIGN KEY (created_by) REFERENCES users(ID) ON DELETE CASCADE
		);`,

		// Conversations table (handles both private and group conversations)
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

		// Group Members table
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

		// Messages table
		`CREATE TABLE IF NOT EXISTS messages (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			conversation_id INTEGER NOT NULL,
			sender_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			message_type VARCHAR(10) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_read BOOLEAN DEFAULT FALSE,
			replied_message_id INTEGER,
			FOREIGN KEY (conversation_id) REFERENCES conversations(ID) ON DELETE CASCADE,
			FOREIGN KEY (sender_id) REFERENCES users(ID) ON DELETE CASCADE,
			FOREIGN KEY (replied_message_id) REFERENCES messages(ID) ON DELETE CASCADE
		);`,

		// Reactions table (for emoji messages)
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
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	return nil
}
