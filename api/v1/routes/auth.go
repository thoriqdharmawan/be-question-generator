package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thoriqdharmawan/be-question-generator/api/v1/controllers"
)

func SetupAuthRoutes(router fiber.Router) {
	router.Post("/login", controllers.Login)
}
