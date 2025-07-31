package chat

import "time"

type SendMessageSchema struct {
	Sender  string `json:"sender" validate:"required,alphanum"`
	Message string `json:"message" validate:"required,min=1,max=250"`
}

type ChatMessage struct {
	ID        uint64 `json:"id"`
	Sender    string `json:"sender" validate:"required,alphanum"`
	Message   string `json:"message" validate:"required,min=1,max=250"`
	CreatedAt string `json:"created_at"`
}

type ChatModel struct {
	ID        uint64 `gorm:"primaryKey"`
	PeerID    string `gorm:"index"`
	Sender    string `gorm:"not null"`
	Message   []byte `gorm:"not null"`
	CreatedAt time.Time
}
