package user

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService *UserService
}

type IUserController interface {
	HandleUserPairing(c *fiber.Ctx) error
}

func NewUserController(userService *UserService) *UserController {
	return &UserController{userService}
}

func (uc *UserController) HandleUserPairing(c *fiber.Ctx) error {
	var input InitPairSchema

	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request pair schema"})
	}

	pubkey, err := uc.userService.HandleUserPairing(input)
	if err != nil {
		switch err.Error() {
		case "pairing code not found":
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "pairing disabled"})
		case "pairing code incorrect":
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		case "user already paired":
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		case "invalid remote public key":
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid public key"})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "unknown error"})
		}
	}

	response := ResponsePairSchema{
		Pubkey: pubkey.Pubkey,
	}

	return c.JSON(response)
}
