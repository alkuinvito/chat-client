package main

import (
	"chat-client/internal/chat"
	"chat-client/internal/discovery"
	"chat-client/internal/router"
	"chat-client/internal/user"
	"chat-client/pkg/db"
	"chat-client/pkg/store"
	"embed"

	"github.com/gofiber/fiber/v2"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Setup new data store
	s := store.NewStore()

	// Setup sqlite db
	db := db.NewDB()

	// Init fiber
	fiberApp := fiber.New(router.DefaultConfig())

	// Init services
	discoveryService := discovery.NewDiscoveryService(s)
	chatService := chat.NewChatService(s, discoveryService)
	userService := user.NewUserService(s, db, fiberApp, discoveryService, chatService)

	// Init controllers
	chatController := chat.NewChatController(chatService)
	userController := user.NewUserController(userService)

	// Init router
	mainRouter := router.NewRouter(fiberApp, chatController, userController)
	mainRouter.Handle()

	// Create an instance of the app structure
	app := NewApp(s, userService, chatService, discoveryService)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "chat-client",
		Width:  900,
		Height: 640,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []any{
			app,
			chatService,
			discoveryService,
			userService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
