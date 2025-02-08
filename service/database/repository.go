package database

import (
	"gorm.io/gorm"
)

type Repository struct {
	database *gorm.DB
}

func NewRepository(database *gorm.DB) *Repository {
	return &Repository{database: database}
}

// toltyramyz
//type AppDatabase interface {
//	//Ping() error
//}
