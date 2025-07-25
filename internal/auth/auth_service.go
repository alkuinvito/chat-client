package auth

import (
	"chat-client/internal/discovery"
	"chat-client/pkg/store"
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type AuthService struct {
	ctx              context.Context
	s                *store.Store
	discoveryService *discovery.DiscoveryService
}

type IAuthService interface {
	Register(username string)
	Startup(ctx context.Context)
}

func NewAuthService(s *store.Store, discoveryService *discovery.DiscoveryService) *AuthService {
	return &AuthService{s: s, discoveryService: discoveryService}
}

func (as *AuthService) Register(username string) {
	as.s.Set("username", username)

	// start broadcasting the service
	go as.discoveryService.BroadcastService(username)

	// redirect to chat page
	runtime.EventsEmit(as.ctx, "navigate", "/chat")
}

func (as *AuthService) Startup(ctx context.Context) {
	as.ctx = ctx
}
