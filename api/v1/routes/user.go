package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thoriqdharmawan/be-question-generator/api/v1/controllers"
	"github.com/thoriqdharmawan/be-question-generator/api/v1/middlewares"
)

func SetupUserRoutes(router fiber.Router) {
	router.Get("/user", controllers.GetUsers)
	router.Post("/user", controllers.CreateUser)
	router.Get("/user/:id", middlewares.Auth, controllers.GetUserById)
	router.Post("/user/verify-email/:token", controllers.VerifyEmailToken)
}
