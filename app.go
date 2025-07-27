package main

import (
	"chat-client/internal/chat"
	"chat-client/internal/discovery"
	"chat-client/internal/user"
	"chat-client/pkg/store"
	"context"
)

// App struct
type App struct {
	ctx              context.Context
	s                *store.Store
	userService      *user.UserService
	chatService      *chat.ChatService
	discoveryService *discovery.DiscoveryService
}

// NewApp creates a new App application struct
func NewApp(s *store.Store, userService *user.UserService, chatService *chat.ChatService, discoveryService *discovery.DiscoveryService) *App {
	return &App{
		s:                s,
		userService:      userService,
		chatService:      chatService,
		discoveryService: discoveryService,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.userService.Startup(ctx)
	a.chatService.Startup(ctx)
	a.discoveryService.Startup(ctx)
}

func (a *App) shutdown(ctx context.Context) {
	a.s.Clear()
}
