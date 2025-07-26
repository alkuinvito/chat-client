package chat

type ChatRoom struct {
	PeerName string `json:"peer_name"`
	IP       string `json:"ip"`
}

type ChatMessage struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}
