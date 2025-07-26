package auth

import (
	"chat-client/internal/chat"
	"chat-client/internal/discovery"
	"chat-client/pkg/store"
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type AuthService struct {
	chatService      *chat.ChatService
	ctx              context.Context
	discoveryService *discovery.DiscoveryService
	s                *store.Store
}

type IAuthService interface {
	Register(username string)
	Startup(ctx context.Context)
}

func NewAuthService(s *store.Store, discoveryService *discovery.DiscoveryService, chatService *chat.ChatService) *AuthService {
	return &AuthService{s: s, discoveryService: discoveryService, chatService: chatService}
}

func (as *AuthService) GetProfile() UserModel {
	var result UserModel

	username, err := as.s.Get("username")
	if err != nil {
		return result
	}

	result.Username = username
	return result

}

func (as *AuthService) Register(username string) {
	as.s.Set("username", username)

	// start broadcasting the service
	go as.discoveryService.BroadcastService(username)

	// start chat server
	go as.chatService.ServeChat()

	// redirect to chat page
	runtime.EventsEmit(as.ctx, "auth:authorized", "/chat")
}

func (as *AuthService) Startup(ctx context.Context) {
	as.ctx = ctx
}
