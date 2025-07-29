package db

import (
	"chat-client/internal/user"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// do auto migrations
	err = db.Find(&user.UserModel{}).Error
	if err != nil {
		log.Println(err)
		db.AutoMigrate(&user.UserModel{})
		db.AutoMigrate(&user.ContactModel{})
	}

	return db
}
