package main

import (
	"chat-client/internal/auth"
	"chat-client/internal/chat"
	"chat-client/internal/discovery"
	"context"
)

// App struct
type App struct {
	ctx              context.Context
	authService      *auth.AuthService
	chatService      *chat.ChatService
	discoveryService *discovery.DiscoveryService
}

// NewApp creates a new App application struct
func NewApp(authService *auth.AuthService, chatService *chat.ChatService, discoveryService *discovery.DiscoveryService) *App {
	return &App{
		authService:      authService,
		chatService:      chatService,
		discoveryService: discoveryService,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.authService.Startup(ctx)
	a.chatService.Startup(ctx)
	a.discoveryService.Startup(ctx)
}
