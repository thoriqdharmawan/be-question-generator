package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thoriqdharmawan/be-question-generator/utils"
)

func Auth(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	if _, err := utils.VerifyJWTToken(token); err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return c.Next()
}
