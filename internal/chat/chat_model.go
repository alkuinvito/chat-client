package chat

import "net"

type ChatRoom struct {
	PeerName string
	IP       net.IP
}

type ChatMessage struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}
