package chat

import (
	"chat-client/internal/discovery"
	"chat-client/pkg/store"
	"context"
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

func (cs *ChatService) GetRooms() []string {
	entries := make(chan *mdns.ServiceEntry, 4)
	var result []string
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for entry := range entries {
			result = append(result, entry.AddrV4.String())
		}
	}()

	go cs.discoveryService.QueryService(entries)

	wg.Wait()

	return result
}

func (cs *ChatService) Startup(ctx context.Context) {
	cs.ctx = ctx
}
