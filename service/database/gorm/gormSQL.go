package gorm

import (
	models2 "WasaText/service/database/models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormSqliteDB() (*gorm.DB, error) {
	dbFile := "database.db" // SQLite This will be your sql file
	//dbFile := "/app/data/database.db" // Matches your Docker volume setup
	Database, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		fmt.Println("Status:", err)
		return nil, err
	}

	err = Database.AutoMigrate(&models2.User{}, &models2.Conversation{}, &models2.Message{}, &models2.GroupMember{}, &models2.Group{}, models2.Reaction{})
	if err != nil {
		return nil, err
	}

	return Database, nil
}
