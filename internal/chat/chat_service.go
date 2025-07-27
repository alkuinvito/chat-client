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
	"strings"
	"sync"

	"github.com/hashicorp/mdns"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ChatService struct {
	ctx              context.Context
	s                *store.Store
	discoveryService *discovery.DiscoveryService
}

type IChatService interface {
	CreateChat(payload ChatMessage)
	GetRooms() response.Response[[]ChatRoom]
	SendMessage(room ChatRoom, message ChatMessage) response.Response[string]
	Startup(ctx context.Context)
}

func NewChatService(s *store.Store, discoveryService *discovery.DiscoveryService) *ChatService {
	return &ChatService{s: s}
}

func (cs *ChatService) CreateChat(payload ChatMessage) {
	runtime.EventsEmit(cs.ctx, "msg:new", fmt.Sprintf(`{"sender":"%s","message":"%s"}`, payload.Sender, payload.Message))
}

func (cs *ChatService) GetRooms() response.Response[[]ChatRoom] {
	entries := make(chan *mdns.ServiceEntry, 4)
	var result []ChatRoom
	var wg sync.WaitGroup

	username, err := cs.s.GetString("username")
	if err != nil {
		return response.New(result).Status(500)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for entry := range entries {
			parsed := strings.Split(entry.Name, ".")
			if len(parsed) == 0 {
				continue
			}

			peerName := parsed[0]
			if peerName == username {
				continue
			}

			result = append(result, ChatRoom{PeerName: peerName, IP: entry.AddrV4.String()})
		}
	}()

	go cs.discoveryService.QueryService(entries)

	wg.Wait()

	return response.New(result)
}

func (cs *ChatService) SendMessage(room ChatRoom, message ChatMessage) response.Response[string] {
	url := fmt.Sprintf("http://%s:%d/api/chat", room.IP, discovery.SVC_PORT)

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
