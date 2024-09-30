package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thoriqdharmawan/be-question-generator/api/v1/controllers"
)

func SetupUserRoutes(router fiber.Router) {
	router.Get("/user", controllers.GetUsers)
}
