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

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/hashicorp/mdns"
)

type ChatService struct {
	ctx              context.Context
	s                *store.Store
	discoveryService *discovery.DiscoveryService
}

type IChatService interface {
	GetRooms() []string
	Startup(ctx context.Context)
}

func NewChatService(s *store.Store, discoveryService *discovery.DiscoveryService) *ChatService {
	return &ChatService{s: s}
}

func (cs *ChatService) GetRooms() []ChatRoom {
	entries := make(chan *mdns.ServiceEntry, 4)
	var result []ChatRoom
	var wg sync.WaitGroup

	username, err := cs.s.Get("username")
	if err != nil {
		return result
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

	return result
}

func (cs *ChatService) SendMessage(room ChatRoom, message ChatMessage) response.Response {
	url := fmt.Sprintf("http://%s:%d/chat", room.IP, discovery.SVC_PORT)

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

func (cs *ChatService) ServeChat() {
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload ChatMessage
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		runtime.EventsEmit(cs.ctx, "msg:new", fmt.Sprintf(`{"sender":"%s","message":"%s"}`, payload.Sender, payload.Message))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"received"}`))
	})

	go func() {
		fmt.Printf("Server is running on: 127.0.0.1:%d\n", discovery.SVC_PORT)
		err := http.ListenAndServe(fmt.Sprintf(":%d", discovery.SVC_PORT), nil)
		if err != nil {
			panic(err)
		}
	}()
}
