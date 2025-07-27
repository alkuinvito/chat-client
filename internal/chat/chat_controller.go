package chat

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ChatController struct {
	chatService *ChatService
}

type IChatController interface {
	CreateChat(c *fiber.Ctx) error
}

func NewChatController(chatService *ChatService) *ChatController {
	return &ChatController{chatService}
}

func (cc *ChatController) CreateChat(c *fiber.Ctx) error {
	var payload ChatMessage

	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid chat message"})
	}

	cc.chatService.CreateChat(payload)

	return c.JSON(fiber.Map{"status": "message received successfully"})
}
