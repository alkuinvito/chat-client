package router

import (
	"chat-client/internal/chat"
	"chat-client/internal/user"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app            *fiber.App
	chatController *chat.ChatController
	userController *user.UserController
}

type IRouter interface {
	Handle()
}

func DefaultConfig() fiber.Config {
	return fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	}
}

func NewRouter(app *fiber.App, chatController *chat.ChatController, userController *user.UserController) *Router {
	return &Router{app, chatController, userController}
}

func (r *Router) Handle() {
	api := r.app.Group("/api")

	chatRouter := api.Group("/chat")
	chatRouter.Post("/send", r.chatController.CreateChat)

	userRouter := api.Group("/user")
	userRouter.Post("/pair", r.userController.PairUser)
}
