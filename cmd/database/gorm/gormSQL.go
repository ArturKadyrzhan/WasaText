package gorm

import (
	"WasaText/cmd/database/models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormSqliteDB() (*gorm.DB, error) {
	dbFile := "database.db" // SQLite This will be your sql file
	Database, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		fmt.Println("Status:", err)
		return nil, err
	}

	err = Database.AutoMigrate(&models.User{}, &models.Conversation{}, &models.Message{}, &models.GroupMember{}, &models.Group{}, models.Reaction{})
	if err != nil {
		return nil, err
	}

	return Database, nil
}
