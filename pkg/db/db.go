package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	return db
}
