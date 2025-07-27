package chat

import (
	"bytes"
	"chat-client/internal/discovery"
	"chat-client/pkg/response"
	"chat-client/pkg/store"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ChatService struct {
	ctx              context.Context
	s                *store.Store
	discoveryService *discovery.DiscoveryService
}

type IChatService interface {
	CreateChat(payload ChatMessage)
	SendMessage(peer discovery.PeerModel, message ChatMessage) response.Response[string]
	Startup(ctx context.Context)
}

func NewChatService(s *store.Store, discoveryService *discovery.DiscoveryService) *ChatService {
	return &ChatService{s: s}
}

func (cs *ChatService) CreateChat(payload ChatMessage) {
	runtime.EventsEmit(cs.ctx, "msg:new", fmt.Sprintf(`{"sender":"%s","message":"%s"}`, payload.Sender, payload.Message))
}

func (cs *ChatService) SendMessage(peer discovery.PeerModel, message ChatMessage) response.Response[string] {
	url := fmt.Sprintf("http://%s:%d/api/chat", peer.IP, discovery.SVC_PORT)

	payload, err := json.Marshal(message)
	if err != nil {
		return response.New("failed to send message").Status(500)
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return response.New("failed to send message").Status(500)
	}

	return response.New("message sent")
}

func (cs *ChatService) Startup(ctx context.Context) {
	cs.ctx = ctx
}
