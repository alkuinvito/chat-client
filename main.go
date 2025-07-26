package main

import (
	"chat-client/internal/auth"
	"chat-client/internal/chat"
	"chat-client/internal/discovery"
	"chat-client/pkg/store"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Setup new data store
	s := store.NewStore()

	// Instantiate services
	discoveryService := discovery.NewDiscoveryService()
	chatService := chat.NewChatService(s, discoveryService)
	authService := auth.NewAuthService(s, discoveryService, chatService)

	// Create an instance of the app structure
	app := NewApp(authService, chatService, discoveryService)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "chat-client",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
			authService,
			chatService,
			discoveryService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
