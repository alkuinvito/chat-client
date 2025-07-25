package main

import (
	"chat-client/internal/auth"
	"chat-client/internal/discovery"
	"context"
)

// App struct
type App struct {
	ctx              context.Context
	authService      *auth.AuthService
	discoveryService *discovery.DiscoveryService
}

// NewApp creates a new App application struct
func NewApp(authService *auth.AuthService, discoveryService *discovery.DiscoveryService) *App {
	return &App{
		authService:      authService,
		discoveryService: discoveryService,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.authService.Startup(ctx)
	a.discoveryService.Startup(ctx)
}
