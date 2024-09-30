package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thoriqdharmawan/be-question-generator/config"
	"github.com/thoriqdharmawan/be-question-generator/utils"
)

func Health(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.Map{
		"ok":  1,
		"v":   config.Conf.Version,
		"env": config.Conf.Environment,
	})
}
