package chat

import (
	"chat-client/internal/discovery"
	"chat-client/pkg/store"
	"context"
	"strings"
	"sync"

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

func (cs *ChatService) Startup(ctx context.Context) {
	cs.ctx = ctx
}
