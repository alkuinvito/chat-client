package chat

import (
	"encoding/json"
	"net/http"
)

type ChatController struct {
	chatService *ChatService
}

type IChatController interface {
	CreateChat() func(w http.ResponseWriter, r *http.Request)
	Routes() http.Handler
}

func NewChatController(chatService *ChatService) *ChatController {
	return &ChatController{chatService}
}

func (cc *ChatController) CreateChat() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload ChatMessage
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		cc.chatService.CreateChat(payload)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"received"}`))
	}
}

func (cc *ChatController) Routes() http.Handler {
	chatRoute := http.NewServeMux()
	chatRoute.HandleFunc("POST /chat", cc.CreateChat())

	return chatRoute
}
