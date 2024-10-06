package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thoriqdharmawan/be-question-generator/api/v1/controllers"
)

func SetupAuthRoutes(router fiber.Router) {
	router.Post("/auth/login", controllers.Login)
	router.Post("/auth/request-reset-password", controllers.RequestResetPassword)
	router.Post("/auth/reset-password", controllers.ResetPassword)
}
