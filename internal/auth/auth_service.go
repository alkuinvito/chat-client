package auth

import (
	"github.com/alkuinvito/chat-client/internal/user"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB, username string) error {
	return db.Create(&user.User{Username: username}).Error
}
