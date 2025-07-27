package user

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService *UserService
}

type IUserController interface {
	PairUser(c *fiber.Ctx) error
}

func NewUserController(userService *UserService) *UserController {
	return &UserController{userService}
}

func (uc *UserController) PairUser(c *fiber.Ctx) error {
	var input RequestPairSchema

	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request pair schema"})
	}

	pubkey, err := uc.userService.PairUser(input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "unknown error"})
	}

	response := ResponsePairSchema{
		Pubkey: pubkey.Pubkey,
	}

	return c.JSON(response)
}
