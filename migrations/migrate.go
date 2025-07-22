package main

import (
	"github.com/alkuinvito/chat-client/internal/user"
	"github.com/alkuinvito/chat-client/pkg/db"
)

func main() {
	db := db.NewDB()
	err := db.AutoMigrate(&user.User{})
	if err != nil {
		panic("migration error")
	}
}
